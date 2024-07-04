package socket

import (
	"context"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Socket struct {
	conn *websocket.Conn
	url  string

	// buffer for the handler func, each received message gets stored into buffer
	buffer chan []byte
	// handler gets executed on each received message from the socket
	handler func([]byte)
	// reconnector gets called when successfully reconnected to the socket
	reconnector func()
}

func NewSocket(
	ctx context.Context,
	url string,
	handler func(bytes []byte),
	reconnector func(),
) (*Socket, error) {
	conn, err := dial(ctx, url)
	if err != nil {
		return nil, err
	}

	socket := &Socket{
		buffer:      make(chan []byte, 1024),
		conn:        conn,
		url:         url,
		handler:     handler,
		reconnector: reconnector,
	}

	go socket.listen(ctx)

	return socket, nil
}

func (w *Socket) SendJSON(ctx context.Context, msg any) error {
	return wsjson.Write(ctx, w.conn, msg)
}

func (w *Socket) reconnect(ctx context.Context) {
	conn, err := dial(ctx, w.url)
	if err != nil {
		time.Sleep(time.Second)
		w.reconnect(ctx)
		return
	}

	w.conn = conn
	w.reconnector()

	go w.listen(ctx)
}

func (w *Socket) readToBuffer(ctx context.Context) {
	for {
		_, b, err := w.conn.Read(ctx)
		if err != nil {
			w.conn.CloseNow()
			if ctx.Err() == nil {
				defer w.reconnect(ctx)
			}
			return
		}
		w.buffer <- b
	}
}

func (w *Socket) listen(ctx context.Context) error {
	go w.readToBuffer(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case bytes := <-w.buffer:
			w.handler(bytes)
		}
	}
}

func dial(ctx context.Context, url string) (*websocket.Conn, error) {
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}
