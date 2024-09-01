package bitvavo

import (
	"fmt"
	"net/url"
	"time"

	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/v2/internal/util"
)

type TradeParams struct {
	// Return the limit most recent trades only.
	// Default: 500
	Limit uint64 `json:"limit"`

	// Return limit trades executed after start.
	Start time.Time `json:"start"`

	// Return limit trades executed before end.
	End time.Time `json:"end"`

	// Return limit trades executed after tradeIdFrom was made.
	TradeIdFrom string `json:"tradeIdFrom"`

	// Return limit trades executed before tradeIdTo was made.
	TradeIdTo string `json:"tradeIdTo"`
}

func (t *TradeParams) Params() url.Values {
	params := make(url.Values)
	if t.Limit > 0 {
		params.Add("limit", fmt.Sprint(t.Limit))
	}
	if !t.Start.IsZero() {
		params.Add("start", fmt.Sprint(t.Start.UnixMilli()))
	}
	if !t.End.IsZero() {
		params.Add("end", fmt.Sprint(t.End.UnixMilli()))
	}
	if t.TradeIdFrom != "" {
		params.Add("tradeIdFrom", t.TradeIdFrom)
	}
	if t.TradeIdTo != "" {
		params.Add("tradeIdTo", t.TradeIdTo)
	}
	return params
}

type TradeHistoric Fill

type Trade struct {
	// The trade ID of the returned trade (UUID).
	Id string `json:"id"`

	// The market which was requested in the subscription.
	Market string `json:"market"`

	// The amount in base currency for which the trade has been made.
	Amount string `json:"amount"`

	// The price in quote currency for which the trade has been made.
	Price string `json:"price"`

	// The side for the taker.
	Side Side `json:"side"`

	// Timestamp in unix milliseconds.
	Timestamp int64 `json:"timestamp"`
}

func (t *Trade) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		id        = util.GetOrEmpty[string]("id", j)
		market    = util.GetOrEmpty[string]("market", j)
		amount    = util.GetOrEmpty[string]("amount", j)
		price     = util.GetOrEmpty[string]("price", j)
		side      = util.GetOrEmpty[string]("side", j)
		timestamp = util.GetOrEmpty[float64]("timestamp", j)
	)

	t.Id = id
	t.Market = market
	t.Amount = amount
	t.Price = price
	t.Side = *sides.Parse(side)
	t.Timestamp = int64(timestamp)

	return nil
}
