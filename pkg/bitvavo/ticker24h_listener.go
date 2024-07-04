package bitvavo

import (
	"context"
	"errors"
	"sync"
)

type Ticker24hEvent ListenerEvent[Ticker24hData]

type Ticker24hListener Listener[Ticker24hEvent]

func NewTicker24hListener() *Ticker24hListener {
	chn := make(chan Ticker24hEvent)
	rchn := make(chan struct{})

	onMessage := func(data WebSocketEventData, err error) {
		if err != nil {
			chn <- Ticker24hEvent{Error: err}
		} else if data.Event == EVENT_TICKER24H {
			var ticker24h Ticker24h
			if err := data.Decode(&ticker24h); err != nil {
				chn <- Ticker24hEvent{Error: err}
			} else {
				for _, t24h := range ticker24h.Data {
					chn <- Ticker24hEvent{Value: t24h}
				}
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws, err := NewWebSocket(ctx, onMessage, func() {
		rchn <- struct{}{}
	})

	if err != nil {
		panic(err)
	}

	return &Ticker24hListener{
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
func (t *Ticker24hListener) Listen(markets []string) (<-chan Ticker24hEvent, error) {
	subs := []Subscription{NewSubscription(CHANNEL_TICKER24H, markets)}
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
func (t *Ticker24hListener) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.subscriptions) == 0 {
		return errors.New("no subscriptions yet, start listening first")
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

func (t *Ticker24hListener) resubscriber() {
	t.once.Do(func() {
		for range t.rchn {
			if err := t.ws.Subscribe(t.subscriptions); err != nil {
				t.chn <- Ticker24hEvent{Error: err}
			}
		}
	})
}
