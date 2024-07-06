package bitvavo

import (
	"context"
	"sync"
)

type CandleEvent ListenerEvent[Candle]

type CandlesListener Listener[CandleEvent]

func NewCandlesListener() *CandlesListener {
	chn := make(chan CandleEvent)
	rchn := make(chan struct{})

	l := &CandlesListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: CHANNEL_CANDLES,
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws, err := NewWebSocket(ctx, l.onMessage, func() {
		rchn <- struct{}{}
	})

	if err != nil {
		panic(err)
	}

	l.closefn = cancel
	l.ws = ws

	return l
}

// Listen for events, you 'can' call this function multiple times.
// The same channel is returned for each function call, meaning that all channel
// receivers get the same data.
func (l *CandlesListener) Listen(markets []string, intervals []Interval) (<-chan CandleEvent, error) {
	if err := l.ws.Subscribe([]Subscription{NewSubscription(l.channel, markets, intervals...)}); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

// Graceful shutdown, once you close a listener it can't be reused, you have to
// create a new one.
func (l *CandlesListener) Close() error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	if err := l.ws.Unsubscribe(l.subscriptions); err != nil {
		return err
	}

	l.closefn()

	close(l.chn)
	close(l.rchn)

	l.subscriptions = nil

	return nil
}

func (l *CandlesListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- CandleEvent{Error: err}
	} else if data.Event == EVENT_SUBSCRIBED {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- CandleEvent{Error: err}
		} else {
			subs, ok := subscribed.SubscriptionsInterval[l.channel]
			if ok {
				markets := make([]string, 0)
				intervals := make([]Interval, 0)
				for i, m := range subs {
					intervals = append(intervals, i)
					markets = append(markets, m...)
				}
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets, intervals...)}
			} else {
				l.chn <- CandleEvent{Error: ErrExpectedChannel(l.channel)}
			}
		}
	} else if data.Event == EVENT_CANDLE {
		var candle Candle
		l.chn <- CandleEvent{Value: candle, Error: data.Decode(&candle)}
	}
}

func (l *CandlesListener) resubscriber() {
	l.once.Do(func() {
		for range l.rchn {
			if err := l.ws.Subscribe(l.subscriptions); err != nil {
				l.chn <- CandleEvent{Error: err}
			}
		}
	})
}
