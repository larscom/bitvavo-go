package bitvavo

import (
	"context"
	"sync"
)

type TradeEvent ListenerEvent[Trade]

type TradesListener Listener[TradeEvent]

func NewTradesListener() *TradesListener {
	chn := make(chan TradeEvent)
	rchn := make(chan struct{})

	l := &TradesListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: CHANNEL_TRADES,
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

// Subscribe to markets, you can call this function multiple times, the same channel is returned.
func (l *TradesListener) Subscribe(markets []string) (<-chan TradeEvent, error) {
	if err := l.ws.Subscribe([]Subscription{NewSubscription(l.channel, markets)}); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

// Unsubscribe from markets.
func (l *TradesListener) Unsubscribe(markets []string) error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	return l.ws.Unsubscribe([]Subscription{NewSubscription(l.channel, markets)})
}

// Graceful shutdown, once you close a listener it can't be reused, you have to
// create a new one.
func (l *TradesListener) Close() error {
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

func (l *TradesListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- TradeEvent{Error: err}
	} else if data.Event == EVENT_SUBSCRIBED || data.Event == EVENT_UNSUBSCRIBED {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- TradeEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.chn <- TradeEvent{Error: ErrExpectedChannel(l.channel)}
			}
		}
	} else if data.Event == EVENT_TRADE {
		var trade Trade
		l.chn <- TradeEvent{Value: trade, Error: data.Decode(&trade)}
	}
}

func (l *TradesListener) resubscriber() {
	l.once.Do(func() {
		for range l.rchn {
			if err := l.ws.Subscribe(l.subscriptions); err != nil {
				l.chn <- TradeEvent{Error: err}
			}
		}
	})
}
