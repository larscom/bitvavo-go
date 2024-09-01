package bitvavo

import (
	"context"
	"sync"
)

type CandleEvent ListenerEvent[Candle]

type CandlesListener listener[CandleEvent]

func NewCandlesListener(options ...WebSocketOption) *CandlesListener {
	chn := make(chan CandleEvent)
	rchn := make(chan struct{})

	l := &CandlesListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: ChannelCandles,
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws, err := NewWebSocket(
		ctx,
		l.onMessage,
		func() { rchn <- struct{}{} },
		options...,
	)

	if err != nil {
		panic(err)
	}

	l.closefn = cancel
	l.ws = ws

	return l
}

// Subscribe to markets with interval.
func (l *CandlesListener) Subscribe(markets []string, intervals []Interval) (<-chan CandleEvent, error) {
	if err := l.ws.Subscribe([]Subscription{NewSubscription(l.channel, markets, intervals...)}); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

// Unsubscribe from markets with intervals.
func (l *CandlesListener) Unsubscribe(markets []string, intervals []Interval) error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	return l.ws.Unsubscribe([]Subscription{NewSubscription(l.channel, markets, intervals...)})
}

// Close everything, graceful shutdown.
func (l *CandlesListener) Close() error {
	defer func() {
		l.closefn()
		close(l.chn)
		close(l.rchn)
	}()

	if len(l.subscriptions) > 0 {
		if err := l.ws.Unsubscribe(l.subscriptions); err != nil {
			return err
		}
	}

	return nil
}

func (l *CandlesListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- CandleEvent{Error: err}
	} else if data.Event == EventSubscribed || data.Event == EventUnsubscribed {
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
				l.subscriptions = nil
			}
		}
	} else if data.Event == EventCandle {
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
