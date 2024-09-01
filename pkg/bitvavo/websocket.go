package bitvavo

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"

	"github.com/larscom/bitvavo-go/v2/internal/crypto"
	"github.com/larscom/bitvavo-go/v2/internal/socket"
	"github.com/larscom/bitvavo-go/v2/internal/util"
	"github.com/orsinium-labs/enum"
)

const websocketURL = "wss://ws.bitvavo.com/v2"

var ErrNotEventType = errors.New("not an event type")

type channelOut struct {
	Name      string   `json:"name"`
	Intervals []string `json:"interval,omitempty"`
	Markets   []string `json:"markets,omitempty"`
}

type messageOut struct {
	Action   string       `json:"action"`
	Channels []channelOut `json:"channels,omitempty"`

	// Api Key.
	Key string `json:"key,omitempty"`
	// SHA256 HMAC hex digest of timestamp + method + url + body.
	Signature string `json:"signature,omitempty"`
	// The current timestamp in milliseconds since 1 Jan 1970.
	Timestamp int64 `json:"timestamp,omitempty"`
}

type WebSocketEvent enum.Member[string]

var (
	webSocketEvent    = enum.NewBuilder[string, WebSocketEvent]()
	EventSubscribed   = webSocketEvent.Add(WebSocketEvent{"subscribed"})
	EventUnsubscribed = webSocketEvent.Add(WebSocketEvent{"unsubscribed"})
	EventCandle       = webSocketEvent.Add(WebSocketEvent{"candle"})
	EventTicker       = webSocketEvent.Add(WebSocketEvent{"ticker"})
	EventTicker24h    = webSocketEvent.Add(WebSocketEvent{"ticker24h"})
	EventTrade        = webSocketEvent.Add(WebSocketEvent{"trade"})
	EventBook         = webSocketEvent.Add(WebSocketEvent{"book"})
	EventAuthenticate = webSocketEvent.Add(WebSocketEvent{"authenticate"})
	EventOrder        = webSocketEvent.Add(WebSocketEvent{"order"})
	EventFill         = webSocketEvent.Add(WebSocketEvent{"fill"})
	webSocketEvents   = webSocketEvent.Enum()
)

type WebSocketEventData struct {
	Event  WebSocketEvent
	Reader io.Reader
}

func (d *WebSocketEventData) Decode(v any) error {
	return json.NewDecoder(d.Reader).Decode(v)
}

func (d *WebSocketEventData) UnmarshalJSON(b []byte) error {
	var j map[string]any

	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	event := util.GetOrEmpty[string]("event", j)
	if event == "" {
		return ErrNotEventType
	}

	d.Event = *webSocketEvents.Parse(event)
	d.Reader = bytes.NewReader(b)

	return nil
}

type WebSocketOption func(*WebSocket)

type WebSocket struct {
	socket     *socket.Socket
	printer    DebugPrinter
	httpClient *http.Client
}

func WithWebSocketHttpClient(client *http.Client) WebSocketOption {
	return func(ws *WebSocket) {
		ws.httpClient = client
	}
}

func WithWebSocketDebugPrinter(printer DebugPrinter) WebSocketOption {
	return func(ws *WebSocket) {
		ws.printer = printer
	}
}

func WithWebSocketDefaultDebugPrinter() WebSocketOption {
	return func(ws *WebSocket) {
		ws.printer = NewDefaultDebugPrinter()
	}
}

func NewWebSocket(
	ctx context.Context,
	messageFunc func(WebSocketEventData, error),
	reconnectFunc func(),
	options ...WebSocketOption,
) (*WebSocket, error) {
	ws := new(WebSocket)
	ws.httpClient = http.DefaultClient

	for _, opt := range options {
		opt(ws)
	}

	onMessage := func(bytes []byte) {
		var data WebSocketEventData
		if err := json.Unmarshal(bytes, &data); err != nil {
			var wsError WebSocketError
			if err := json.Unmarshal(bytes, &wsError); err != nil {
				messageFunc(data, err)
			} else {
				messageFunc(data, &wsError)
			}
		} else {
			messageFunc(data, nil)
		}
	}

	onDebug := func(message string) {
		debug(ws.printer, message)
	}

	opts := &socket.Options{
		Url:           websocketURL,
		HttpClient:    ws.httpClient,
		MessageFunc:   onMessage,
		ReconnectFunc: reconnectFunc,
		DebugFunc:     onDebug,
	}
	s, err := socket.NewSocket(ctx, opts)
	if err != nil {
		return nil, err
	}
	ws.socket = s

	return ws, nil
}

func (w *WebSocket) Authenticate(apiKey string, apiSecret string) error {
	timestamp := time.Now().UnixMilli()
	msg := messageOut{
		Action:    "authenticate",
		Key:       apiKey,
		Signature: crypto.CreateSignature("GET", "/websocket", nil, timestamp, apiSecret),
		Timestamp: timestamp,
	}

	return w.socket.SendJSON(context.Background(), msg)
}

func (w *WebSocket) Subscribe(subscriptions []Subscription) error {
	msg := messageOut{
		Action:   "subscribe",
		Channels: mapToChannels(subscriptions),
	}

	return w.socket.SendJSON(context.Background(), msg)
}

func (w *WebSocket) Unsubscribe(subscriptions []Subscription) error {
	msg := messageOut{
		Action:   "unsubscribe",
		Channels: mapToChannels(subscriptions),
	}

	return w.socket.SendJSON(context.Background(), msg)
}

func mapToChannels(subscriptions []Subscription) []channelOut {
	channels := make([]channelOut, len(subscriptions))

	for i, subscription := range subscriptions {
		intervals := make([]string, len(subscription.Intervals))
		for y, interval := range subscription.Intervals {
			intervals[y] = interval.Value
		}
		channels[i] = channelOut{
			Name:      subscription.Channel.Value,
			Markets:   subscription.Markets,
			Intervals: intervals,
		}
	}

	return channels
}
