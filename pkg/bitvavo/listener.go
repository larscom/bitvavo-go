package bitvavo

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNoSubscriptions = errors.New("no subscriptions yet, start listening first")
	ErrNoAuth          = errors.New("received auth event from server, but was not authenticated")
	ErrExpectedChannel = func(chn Channel) error { return fmt.Errorf("expected channel '%s' in subscribed event", chn.Value) }
)

type Listener[T any] interface {
	Subscriber[T]
	Unsubscriber
	Closer
}

type Subscriber[T any] interface {
	// Subscribe to markets
	Subscribe(markets []string) (<-chan T, error)
}

type Unsubscriber interface {
	// Unsubscribe from markets
	Unsubscribe(markets []string) error
}

type Closer interface {
	// Close everything, graceful shutdown.
	Close() error
}

type ListenerEvent[T any] struct {
	Value T
	Error error
}

type listener[T any] struct {
	ws            *WebSocket
	chn           chan T
	rchn          chan struct{}
	once          *sync.Once
	channel       Channel
	subscriptions []Subscription
	closefn       context.CancelFunc
}

type authListener[T any] struct {
	listener[T]
	apiKey      string
	apiSecret   string
	authchn     chan bool
	pendingsubs chan []Subscription
}
