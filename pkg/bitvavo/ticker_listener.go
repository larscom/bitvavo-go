package bitvavo

import (
	"context"
	"sync"
)

type TickerEvent ListenerEvent[Ticker]

type TickerListener Listener[TickerEvent]

func NewTickerListener() *TickerListener {
	chn := make(chan TickerEvent)
	rchn := make(chan struct{})

	l := &TickerListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: CHANNEL_TICKER,
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
func (l *TickerListener) Listen(markets []string) (<-chan TickerEvent, error) {
	if err := l.ws.Subscribe([]Subscription{NewSubscription(l.channel, markets)}); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

// Graceful shutdown, once you close a listener it can't be reused, you have to
// create a new one.
func (l *TickerListener) Close() error {
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

func (l *TickerListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- TickerEvent{Error: err}
	} else if data.Event == EVENT_SUBSCRIBED {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- TickerEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.chn <- TickerEvent{Error: ErrExpectedChannel(l.channel)}
			}
		}
	} else if data.Event == EVENT_TICKER {
		var ticker Ticker
		l.chn <- TickerEvent{Value: ticker, Error: data.Decode(&ticker)}
	}
}

func (t *TickerListener) resubscriber() {
	t.once.Do(func() {
		for range t.rchn {
			if err := t.ws.Subscribe(t.subscriptions); err != nil {
				t.chn <- TickerEvent{Error: err}
			}
		}
	})
}
