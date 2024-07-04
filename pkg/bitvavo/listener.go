package bitvavo

import (
	"context"
	"sync"
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
