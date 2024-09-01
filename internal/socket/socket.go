package socket

import (
	"context"
	"errors"
	"fmt"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"net/http"
	"time"
)

const readLimit = 655350

type Options struct {
	Url        string
	HttpClient *http.Client

	// messageFunc gets executed on each received message from the socket
	MessageFunc func(bytes []byte)
	// reconnectFunc gets called when successfully reconnected to the socket
	ReconnectFunc func()
	// debugFunc gets called on every connection event
	DebugFunc func(string)
}

type Socket struct {
	conn *websocket.Conn
	// buffer for the messageFunc func, each received message gets stored into buffer
	buffer chan []byte

	options *Options
}

func NewSocket(
	ctx context.Context,
	options *Options,
) (*Socket, error) {
	if options == nil {
		return nil, errors.New("options is nil")
	}

	conn, err := dial(ctx, options.Url, options.HttpClient)
	if err != nil {
		return nil, err
	}

	socket := &Socket{
		buffer:  make(chan []byte, 1024),
		conn:    conn,
		options: options,
	}

	go func() {
		_ = socket.listen(ctx)
	}()

	return socket, nil
}

func (w *Socket) SendJSON(ctx context.Context, msg any) error {
	return wsjson.Write(ctx, w.conn, msg)
}

func (w *Socket) reconnect(ctx context.Context) {
	w.options.DebugFunc("websocket reconnecting...")

	conn, err := dial(ctx, w.options.Url, w.options.HttpClient)
	if err != nil {
		w.options.DebugFunc(fmt.Sprint("websocket error while reconnecting: ", err))
		time.Sleep(time.Second)
		w.reconnect(ctx)
		return
	}

	w.conn = conn
	w.options.ReconnectFunc()
	w.options.DebugFunc("websocket reconnected")

	go func() {
		_ = w.listen(ctx)
	}()
}

func (w *Socket) readToBuffer(ctx context.Context) {
	w.options.DebugFunc("websocket connected")
	for {
		_, b, err := w.conn.Read(ctx)
		if err != nil {
			_ = w.conn.CloseNow()
			w.options.DebugFunc(fmt.Sprint("websocket disconnected with error: ", err))
			if ctx.Err() == nil {
				//goland:noinspection ALL
				defer w.reconnect(ctx)
			}
			return
		}
		w.buffer <- b
		w.options.DebugFunc(fmt.Sprint("websocket received: ", string(b)))
	}
}

func (w *Socket) listen(ctx context.Context) error {
	go w.readToBuffer(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case bytes := <-w.buffer:
			w.options.MessageFunc(bytes)
		}
	}
}

func dial(ctx context.Context, url string, httpClient *http.Client) (*websocket.Conn, error) {
	c, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, err
	}
	c.SetReadLimit(readLimit)
	return c, nil
}
