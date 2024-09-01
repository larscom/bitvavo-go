package bitvavo

import (
	"context"
	"sync"
)

type Ticker24hEvent ListenerEvent[Ticker24hData]

type Ticker24hListener listener[Ticker24hEvent]

func NewTicker24hListener(options ...WebSocketOption) Listener[Ticker24hEvent] {
	chn := make(chan Ticker24hEvent)
	rchn := make(chan struct{})

	l := &Ticker24hListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: ChannelTicker24h,
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

func (l *Ticker24hListener) Subscribe(markets []string) (<-chan Ticker24hEvent, error) {
	subs := []Subscription{NewSubscription(l.channel, markets)}
	if err := l.ws.Subscribe(subs); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

func (l *Ticker24hListener) Unsubscribe(markets []string) error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	return l.ws.Unsubscribe([]Subscription{NewSubscription(l.channel, markets)})
}

func (l *Ticker24hListener) Close() error {
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

func (l *Ticker24hListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- Ticker24hEvent{Error: err}
	} else if data.Event == EventSubscribed || data.Event == EventUnsubscribed {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- Ticker24hEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.subscriptions = nil
			}
		}
	} else if data.Event == EventTicker24h {
		var ticker24h Ticker24h
		if err := data.Decode(&ticker24h); err != nil {
			l.chn <- Ticker24hEvent{Error: err}
		} else {
			for _, t24h := range ticker24h.Data {
				l.chn <- Ticker24hEvent{Value: t24h}
			}
		}
	}
}

func (l *Ticker24hListener) resubscriber() {
	l.once.Do(func() {
		for range l.rchn {
			if err := l.ws.Subscribe(l.subscriptions); err != nil {
				l.chn <- Ticker24hEvent{Error: err}
			}
		}
	})
}
