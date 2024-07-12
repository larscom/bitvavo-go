package bitvavo

import (
	"context"
	"sync"
)

type FillEvent ListenerEvent[Fill]

type FillListener AuthListener[FillEvent]

func NewFillListener(apiKey, apiSecret string) *FillListener {
	chn := make(chan FillEvent)
	rchn := make(chan struct{})
	authchn := make(chan bool)
	pendingsubs := make(chan []Subscription)

	l := &FillListener{
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		authchn:     authchn,
		pendingsubs: pendingsubs,
		Listener: Listener[FillEvent]{
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

// Subscribe to markets, you can call this function multiple times, the same channel is returned.
func (l *FillListener) Subscribe(markets []string) (<-chan FillEvent, error) {
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

// Unsubscribe from markets.
func (l *FillListener) Unsubscribe(markets []string) error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	return l.ws.Unsubscribe([]Subscription{NewSubscription(l.channel, markets)})
}

// Graceful shutdown, once you close a listener it can't be reused, you have to
// create a new one.
func (l *FillListener) Close() error {
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

func (l *FillListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- FillEvent{Error: err}
	} else if data.Event == EVENT_AUTHENTICATE {
		var auth Authenticate
		if err := data.Decode(&auth); err != nil {
			l.chn <- FillEvent{Error: err}
		} else {
			l.authchn <- auth.Authenticated
		}
	} else if data.Event == EVENT_SUBSCRIBED || data.Event == EVENT_UNSUBSCRIBED {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- FillEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.chn <- FillEvent{Error: ErrExpectedChannel(l.channel)}
			}
		}
	} else if data.Event == EVENT_FILL {
		var fill Fill
		l.chn <- FillEvent{Value: fill, Error: data.Decode(&fill)}
	}
}

// First authenticate on reconnect, then we receive a authenticated event from the server.
// If we are successfully authenticated we do a subscribe.
func (l *FillListener) resubscriber() {
	l.once.Do(func() {
		for {
			select {
			case <-l.rchn:
				if err := l.ws.Authenticate(l.apiKey, l.apiSecret); err != nil {
					l.chn <- FillEvent{Error: err}
				} else {
					l.pendingsubs <- l.subscriptions
				}
			case authenticated := <-l.authchn:
				pendingSubs := <-l.pendingsubs
				if authenticated {
					if err := l.ws.Subscribe(pendingSubs); err != nil {
						l.chn <- FillEvent{Error: err}
					}
				} else {
					l.chn <- FillEvent{Error: ErrNoAuth}
				}
			}
		}
	})
}
