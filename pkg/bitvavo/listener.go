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

type ListenerEvent[T any] struct {
	Value T
	Error error
}

type Listener[T any] struct {
	ws            *WebSocket
	chn           chan (T)
	rchn          chan (struct{})
	once          *sync.Once
	mu            sync.RWMutex
	subscriptions []Subscription
	closefn       context.CancelFunc
}

type AuthListener[T any] struct {
	Listener[T]
	apiKey    string
	apiSecret string
	authchn   chan (bool)
}
