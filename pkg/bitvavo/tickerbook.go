package bitvavo

import (
	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/util"
)

type TickerBook struct {
	// The market you requested the current best orders for.
	Market string `json:"market"`

	// The highest buy order in quote currency for market currently available on Bitvavo.
	Bid float64 `json:"bid"`

	// The amount of base currency for bid in the order.
	BidSize float64 `json:"bidSize"`

	// The lowest sell order in quote currency for market currently available on Bitvavo.
	Ask float64 `json:"ask"`

	// The amount of base currency for ask in the order.
	AskSize float64 `json:"askSize"`
}

func (t *TickerBook) UnmarshalJSON(bytes []byte) error {
	var j map[string]string

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		market  = j["market"]
		bid     = j["bid"]
		bidSize = j["bidSize"]
		ask     = j["ask"]
		askSize = j["askSize"]
	)

	t.Market = market
	t.Bid = util.IfOrElse(len(bid) > 0, func() float64 { return util.MustFloat64(bid) }, 0)
	t.BidSize = util.IfOrElse(len(bidSize) > 0, func() float64 { return util.MustFloat64(bidSize) }, 0)
	t.Ask = util.IfOrElse(len(ask) > 0, func() float64 { return util.MustFloat64(ask) }, 0)
	t.AskSize = util.IfOrElse(len(askSize) > 0, func() float64 { return util.MustFloat64(askSize) }, 0)

	return nil
}
