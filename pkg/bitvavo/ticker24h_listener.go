package bitvavo

import (
	"context"
	"sync"
)

type Ticker24hEvent ListenerEvent[Ticker24hData]

type Ticker24hListener Listener[Ticker24hEvent]

func NewTicker24hListener() *Ticker24hListener {
	chn := make(chan Ticker24hEvent)
	rchn := make(chan struct{})

	l := &Ticker24hListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: CHANNEL_TICKER24H,
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
func (l *Ticker24hListener) Listen(markets []string) (<-chan Ticker24hEvent, error) {
	subs := []Subscription{NewSubscription(l.channel, markets)}
	if err := l.ws.Subscribe(subs); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

// Graceful shutdown, once you close a listener it can't be reused, you have to
// create a new one.
func (l *Ticker24hListener) Close() error {
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

func (l *Ticker24hListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- Ticker24hEvent{Error: err}
	} else if data.Event == EVENT_SUBSCRIBED {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- Ticker24hEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.chn <- Ticker24hEvent{Error: ErrExpectedChannel(l.channel)}
			}
		}
	} else if data.Event == EVENT_TICKER24H {
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
