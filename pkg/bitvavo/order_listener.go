package bitvavo

import (
	"context"
	"sync"
)

type OrderEvent ListenerEvent[Order]

type OrderListener authListener[OrderEvent]

func NewOrderListener(apiKey, apiSecret string) Listener[OrderEvent] {
	chn := make(chan OrderEvent)
	rchn := make(chan struct{})
	authchn := make(chan bool)
	pendingsubs := make(chan []Subscription)

	l := &OrderListener{
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		authchn:     authchn,
		pendingsubs: pendingsubs,
		listener: listener[OrderEvent]{
			chn:     chn,
			rchn:    rchn,
			once:    new(sync.Once),
			channel: CHANNEL_ACCOUNT,
		},
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

func (l *OrderListener) Subscribe(markets []string) (<-chan OrderEvent, error) {
	if err := l.ws.Authenticate(l.apiKey, l.apiSecret); err != nil {
		return nil, err
	}

	go func() {
		// blocks until we receive an authenticated event
		l.pendingsubs <- []Subscription{NewSubscription(l.channel, markets)}
	}()

	go l.resubscriber()

	return l.chn, nil
}

func (l *OrderListener) Unsubscribe(markets []string) error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	return l.ws.Unsubscribe([]Subscription{NewSubscription(l.channel, markets)})
}

func (l *OrderListener) Close() error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	if err := l.ws.Unsubscribe(l.subscriptions); err != nil {
		return err
	}

	l.closefn()

	close(l.chn)
	close(l.rchn)
	close(l.authchn)
	close(l.pendingsubs)

	l.subscriptions = nil

	return nil
}

func (l *OrderListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- OrderEvent{Error: err}
	} else if data.Event == EVENT_AUTHENTICATE {
		var auth Authenticate
		if err := data.Decode(&auth); err != nil {
			l.chn <- OrderEvent{Error: err}
		} else {
			l.authchn <- auth.Authenticated
		}
	} else if data.Event == EVENT_SUBSCRIBED || data.Event == EVENT_UNSUBSCRIBED {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- OrderEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.chn <- OrderEvent{Error: ErrExpectedChannel(l.channel)}
			}
		}
	} else if data.Event == EVENT_ORDER {
		var order Order
		l.chn <- OrderEvent{Value: order, Error: data.Decode(&order)}
	}
}

// First authenticate on reconnect, then we receive a authenticated event from the server.
// If we are successfully authenticated we do a subscribe.
func (l *OrderListener) resubscriber() {
	l.once.Do(func() {
		for {
			select {
			case <-l.rchn:
				if err := l.ws.Authenticate(l.apiKey, l.apiSecret); err != nil {
					l.chn <- OrderEvent{Error: err}
				} else {
					l.pendingsubs <- l.subscriptions
				}
			case authenticated := <-l.authchn:
				pendingSubs := <-l.pendingsubs
				if authenticated {
					if err := l.ws.Subscribe(pendingSubs); err != nil {
						l.chn <- OrderEvent{Error: err}
					}
				} else {
					l.chn <- OrderEvent{Error: ErrNoAuth}
				}
			}
		}
	})
}
