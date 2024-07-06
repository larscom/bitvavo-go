package bitvavo

import (
	"context"
	"sync"
)

type OrderEvent ListenerEvent[Order]

type OrderListener AuthListener[OrderEvent]

func NewOrderListener(apiKey, apiSecret string) *OrderListener {
	chn := make(chan OrderEvent)
	rchn := make(chan struct{})
	authchn := make(chan bool)

	onMessage := func(data WebSocketEventData, err error) {
		if err != nil {
			chn <- OrderEvent{Error: err}
		} else if data.Event == EVENT_AUTHENTICATE {
			var auth Authenticate
			if err := data.Decode(&auth); err != nil {
				chn <- OrderEvent{Error: err}
			} else {
				authchn <- auth.Authenticated
			}
		} else if data.Event == EVENT_ORDER {
			var order Order
			chn <- OrderEvent{Value: order, Error: data.Decode(&order)}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws, err := NewWebSocket(ctx, onMessage, func() {
		rchn <- struct{}{}
	})

	if err != nil {
		panic(err)
	}

	return &OrderListener{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		authchn:   authchn,
		Listener: Listener[OrderEvent]{
			ws:            ws,
			chn:           chn,
			rchn:          rchn,
			once:          &sync.Once{},
			subscriptions: make([]Subscription, 0),
			closefn:       cancel,
		},
	}
}

// Listen for events, you 'can' call this function multiple times.
// The same channel is returned for each function call, meaning that all channel
// receivers get the same data.
func (t *OrderListener) Listen(markets []string) (<-chan OrderEvent, error) {
	if err := t.ws.Authenticate(t.apiKey, t.apiSecret); err != nil {
		return nil, err
	}

	subs := []Subscription{NewSubscription(CHANNEL_ACCOUNT, markets)}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.subscriptions = append(t.subscriptions, subs...)

	go t.resubscriber()

	return t.chn, nil
}

// Graceful shutdown, once you close a listener it can't be reused, you have to
// create a new one.
func (t *OrderListener) Close() error {
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
	close(t.authchn)

	t.subscriptions = nil

	return nil
}

// First authenticate on reconnect, then we receive a authenticated event from the server.
// If we are successfully authenticated we do a subscribe.
func (t *OrderListener) resubscriber() {
	t.once.Do(func() {
		for {
			select {
			case <-t.rchn:
				if err := t.ws.Authenticate(t.apiKey, t.apiSecret); err != nil {
					t.chn <- OrderEvent{Error: err}
				}
			case authenticated := <-t.authchn:
				if authenticated {
					if err := t.ws.Subscribe(t.subscriptions); err != nil {
						t.chn <- OrderEvent{Error: err}
					}
				} else {
					t.chn <- OrderEvent{Error: ErrNoAuth}
				}
			}
		}
	})
}
