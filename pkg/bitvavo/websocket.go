package bitvavo

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/goccy/go-json"

	"github.com/larscom/bitvavo-go/internal/crypto"
	"github.com/larscom/bitvavo-go/internal/socket"
	"github.com/larscom/bitvavo-go/internal/util"
	"github.com/orsinium-labs/enum"
)

const bitvavoWSURL = "wss://ws.bitvavo.com/v2"

type chn struct {
	Name      string   `json:"name"`
	Intervals []string `json:"interval,omitempty"`
	Markets   []string `json:"markets,omitempty"`
}

type message struct {
	Action   string `json:"action"`
	Channels []chn  `json:"channels,omitempty"`

	// Api Key.
	Key string `json:"key,omitempty"`
	// SHA256 HMAC hex digest of timestamp + method + url + body.
	Signature string `json:"signature,omitempty"`
	// The current timestamp in milliseconds since 1 Jan 1970.
	Timestamp int64 `json:"timestamp,omitempty"`
}

type WebSocketEvent enum.Member[string]

var (
	webSocketEvent     = enum.NewBuilder[string, WebSocketEvent]()
	EVENT_SUBSCRIBED   = webSocketEvent.Add(WebSocketEvent{"subscribed"})
	EVENT_UNSUBSCRIBED = webSocketEvent.Add(WebSocketEvent{"unsubscribed"})
	EVENT_CANDLE       = webSocketEvent.Add(WebSocketEvent{"candle"})
	EVENT_TICKER       = webSocketEvent.Add(WebSocketEvent{"ticker"})
	EVENT_TICKER24H    = webSocketEvent.Add(WebSocketEvent{"ticker24h"})
	EVENT_TRADE        = webSocketEvent.Add(WebSocketEvent{"trade"})
	EVENT_BOOK         = webSocketEvent.Add(WebSocketEvent{"book"})
	EVENT_AUTHENTICATE = webSocketEvent.Add(WebSocketEvent{"authenticate"})
	EVENT_ORDER        = webSocketEvent.Add(WebSocketEvent{"order"})
	EVENT_FILL         = webSocketEvent.Add(WebSocketEvent{"fill"})
	webSocketEvents    = webSocketEvent.Enum()
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
		return errors.New("not an event type")
	}

	d.Event = *webSocketEvents.Parse(event)
	d.Reader = bytes.NewReader(b)

	return nil
}

type WebSocket struct {
	socket *socket.Socket
}

func NewWebSocket(ctx context.Context, sendMessage func(WebSocketEventData, error), reconnFunc func()) (*WebSocket, error) {
	onMessage := func(bytes []byte) {
		var data WebSocketEventData
		if err := json.Unmarshal(bytes, &data); err != nil {
			var wsError WebSocketError
			if err := json.Unmarshal(bytes, &wsError); err != nil {
				sendMessage(data, err)
			} else {
				sendMessage(data, &wsError)
			}
		} else {
			sendMessage(data, nil)
		}
	}

	socket, err := socket.NewSocket(ctx, bitvavoWSURL, onMessage, reconnFunc)
	if err != nil {
		return nil, err
	}

	return &WebSocket{
		socket: socket,
	}, nil
}

func (w *WebSocket) Authenticate(apiKey string, apiSecret string) error {
	timestamp := time.Now().UnixMilli()
	msg := message{
		Action:    "authenticate",
		Key:       apiKey,
		Signature: crypto.CreateSignature("GET", "/websocket", nil, timestamp, apiSecret),
		Timestamp: timestamp,
	}

	return w.socket.SendJSON(context.Background(), msg)
}

func (w *WebSocket) Subscribe(subscriptions []Subscription) error {
	msg := message{
		Action:   "subscribe",
		Channels: mapToChannels(subscriptions),
	}

	return w.socket.SendJSON(context.Background(), msg)
}

func (w *WebSocket) Unsubscribe(subscriptions []Subscription) error {
	msg := message{
		Action:   "unsubscribe",
		Channels: mapToChannels(subscriptions),
	}

	return w.socket.SendJSON(context.Background(), msg)
}

func mapToChannels(subscriptions []Subscription) []chn {
	channels := make([]chn, len(subscriptions))

	for i, subscription := range subscriptions {
		intervals := make([]string, len(subscription.Intervals))
		for y, interval := range subscription.Intervals {
			intervals[y] = interval.Value
		}
		channels[i] = chn{
			Name:      subscription.Channel.Value,
			Markets:   subscription.Markets,
			Intervals: intervals,
		}
	}

	return channels
}
