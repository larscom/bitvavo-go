package bitvavo

import (
	"context"
	"sync"
)

type BookEvent ListenerEvent[Book]

type BookListener listener[BookEvent]

func NewBookListener() Listener[BookEvent] {
	chn := make(chan BookEvent)
	rchn := make(chan struct{})

	l := &BookListener{
		chn:     chn,
		rchn:    rchn,
		once:    new(sync.Once),
		channel: ChannelBook,
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

func (l *BookListener) Subscribe(markets []string) (<-chan BookEvent, error) {
	if err := l.ws.Subscribe([]Subscription{NewSubscription(l.channel, markets)}); err != nil {
		return nil, err
	}

	go l.resubscriber()

	return l.chn, nil
}

func (l *BookListener) Unsubscribe(markets []string) error {
	if len(l.subscriptions) == 0 {
		return ErrNoSubscriptions
	}

	return l.ws.Unsubscribe([]Subscription{NewSubscription(l.channel, markets)})
}

func (l *BookListener) Close() error {
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

func (l *BookListener) onMessage(data WebSocketEventData, err error) {
	if err != nil {
		l.chn <- BookEvent{Error: err}
	} else if data.Event == EventSubscribed || data.Event == EventUnsubscribed {
		var subscribed Subscribed
		if err := data.Decode(&subscribed); err != nil {
			l.chn <- BookEvent{Error: err}
		} else {
			markets, ok := subscribed.Subscriptions[l.channel]
			if ok {
				l.subscriptions = []Subscription{NewSubscription(l.channel, markets)}
			} else {
				l.chn <- BookEvent{Error: ErrExpectedChannel(l.channel)}
			}
		}
	} else if data.Event == EventBook {
		var book Book
		l.chn <- BookEvent{Value: book, Error: data.Decode(&book)}
	}
}

func (l *BookListener) resubscriber() {
	l.once.Do(func() {
		for range l.rchn {
			if err := l.ws.Subscribe(l.subscriptions); err != nil {
				l.chn <- BookEvent{Error: err}
			}
		}
	})
}
