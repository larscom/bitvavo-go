package bitvavo

import (
	"context"
	"sync"
)

type TickerEvent ListenerEvent[Ticker]

type TickerListener listener[TickerEvent]

func NewTickerListener(options ...WebSocketOption) Listener[TickerEvent] {
	chn := make(chan TickerEvent)
	rchn := make(chan struct{})

	l := &TickerListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: ChannelTicker,
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

func (l *TickerListener) Subscribe(markets []string) (<-chan TickerEvent, error) {
	if err := l.ws.Subscribe([]Subscription{NewSubscription(l.channel, markets)}); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

func (l *TickerListener) Unsubscribe(markets []string) error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	return l.ws.Unsubscribe([]Subscription{NewSubscription(l.channel, markets)})
}

func (l *TickerListener) Close() error {
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

func (l *TickerListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- TickerEvent{Error: err}
	} else if data.Event == EventSubscribed || data.Event == EventUnsubscribed {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- TickerEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.subscriptions = nil
			}
		}
	} else if data.Event == EventTicker {
		var ticker Ticker
		l.chn <- TickerEvent{Value: ticker, Error: data.Decode(&ticker)}
	}
}

func (l *TickerListener) resubscriber() {
	l.once.Do(func() {
		for range l.rchn {
			if err := l.ws.Subscribe(l.subscriptions); err != nil {
				l.chn <- TickerEvent{Error: err}
			}
		}
	})
}
