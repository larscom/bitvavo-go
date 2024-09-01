package bitvavo

import (
	"context"
	"sync"
)

type TradeEvent ListenerEvent[Trade]

type TradesListener listener[TradeEvent]

func NewTradesListener(options ...WebSocketOption) Listener[TradeEvent] {
	chn := make(chan TradeEvent)
	rchn := make(chan struct{})

	l := &TradesListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: ChannelTrades,
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

func (l *TradesListener) Subscribe(markets []string) (<-chan TradeEvent, error) {
	if err := l.ws.Subscribe([]Subscription{NewSubscription(l.channel, markets)}); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

func (l *TradesListener) Unsubscribe(markets []string) error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	return l.ws.Unsubscribe([]Subscription{NewSubscription(l.channel, markets)})
}

func (l *TradesListener) Close() error {
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

func (l *TradesListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- TradeEvent{Error: err}
	} else if data.Event == EventSubscribed || data.Event == EventUnsubscribed {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- TradeEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.subscriptions = nil
			}
		}
	} else if data.Event == EventTrade {
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
