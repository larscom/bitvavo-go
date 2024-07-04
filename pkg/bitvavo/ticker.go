package bitvavo

import (
	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/util"
)

type Ticker struct {
	// The market which was requested in the subscription.
	Market string `json:"market"`

	// The price of the best (highest) bid offer available, only sent when either bestBid or bestBidSize has changed.
	BestBid float64 `json:"bestBid"`

	// The size of the best (highest) bid offer available, only sent when either bestBid or bestBidSize has changed.
	BestBidSize float64 `json:"bestBidSize"`

	// The price of the best (lowest) ask offer available, only sent when either bestAsk or bestAskSize has changed.
	BestAsk float64 `json:"bestAsk"`

	// The size of the best (lowest) ask offer available, only sent when either bestAsk or bestAskSize has changed.
	BestAskSize float64 `json:"bestAskSize"`

	// The last price for which a trade has occurred, only sent when lastPrice has changed.
	LastPrice float64 `json:"lastPrice"`
}

func (t *Ticker) UnmarshalJSON(bytes []byte) error {
	var j map[string]string

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		market      = j["market"]
		bestBid     = j["bestBid"]
		bestBidSize = j["bestBidSize"]
		bestAsk     = j["bestAsk"]
		bestAskSize = j["bestAskSize"]
		lastPrice   = j["lastPrice"]
	)

	t.Market = market
	t.BestBid = util.IfOrElse(len(bestBid) > 0, func() float64 { return util.MustFloat64(bestBid) }, 0)
	t.BestBidSize = util.IfOrElse(len(bestBidSize) > 0, func() float64 { return util.MustFloat64(bestBidSize) }, 0)
	t.BestAsk = util.IfOrElse(len(bestAsk) > 0, func() float64 { return util.MustFloat64(bestAsk) }, 0)
	t.BestAskSize = util.IfOrElse(len(bestAskSize) > 0, func() float64 { return util.MustFloat64(bestAskSize) }, 0)
	t.LastPrice = util.IfOrElse(len(lastPrice) > 0, func() float64 { return util.MustFloat64(lastPrice) }, 0)

	return nil
}
