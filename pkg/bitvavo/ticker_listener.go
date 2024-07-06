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

	onMessage := func(data WebSocketEventData, err error) {
		if err != nil {
			chn <- TickerEvent{Error: err}
		} else if data.Event == EVENT_TICKER {
			var ticker Ticker
			chn <- TickerEvent{Value: ticker, Error: data.Decode(&ticker)}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws, err := NewWebSocket(ctx, onMessage, func() {
		rchn <- struct{}{}
	})

	if err != nil {
		panic(err)
	}

	return &TickerListener{
		ws:            ws,
		chn:           chn,
		rchn:          rchn,
		once:          &sync.Once{},
		subscriptions: make([]Subscription, 0),
		closefn:       cancel,
	}
}

// Listen for events, you 'can' call this function multiple times.
// The same channel is returned for each function call, meaning that all channel
// receivers get the same data.
func (t *TickerListener) Listen(markets []string) (<-chan TickerEvent, error) {
	subs := []Subscription{NewSubscription(CHANNEL_TICKER, markets)}
	if err := t.ws.Subscribe(subs); err != nil {
		return nil, err
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.subscriptions = append(t.subscriptions, subs...)

	go t.resubscriber()

	return t.chn, nil
}

// Graceful shutdown, once you close a listener it can't be reused, you have to
// create a new one.
func (t *TickerListener) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	if err := t.ws.Unsubscribe(t.subscriptions); err != nil {
		return err
	}

	t.closefn()

	close(t.chn)
	close(t.rchn)

	t.subscriptions = nil

	return nil
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
