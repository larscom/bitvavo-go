package bitvavo

import (
	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/util"
)

type Ticker24h struct {
	Data []Ticker24hData `json:"data"`
}

type Ticker24hData struct {
	// The market which was requested in the subscription.
	Market string `json:"market"`

	// The open price of the 24-hour period.
	Open float64 `json:"open"`

	// The highest price for which a trade occurred in the 24-hour period.
	High float64 `json:"high"`

	// The lowest price for which a trade occurred in the 24-hour period.
	Low float64 `json:"low"`

	// The last price for which a trade occurred in the 24-hour period.
	Last float64 `json:"last"`

	// The total volume of the 24-hour period in base currency.
	Volume float64 `json:"volume"`

	// The total volume of the 24-hour period in quote currency.
	VolumeQuote float64 `json:"volumeQuote"`

	// The best (highest) bid offer at the current moment.
	Bid float64 `json:"bid"`

	// The size of the best (highest) bid offer.
	BidSize float64 `json:"bidSize"`

	// The best (lowest) ask offer at the current moment.
	Ask float64 `json:"ask"`

	// The size of the best (lowest) ask offer.
	AskSize float64 `json:"askSize"`

	// Timestamp in unix milliseconds.
	Timestamp int64 `json:"timestamp"`

	// Start timestamp in unix milliseconds.
	StartTimestamp int64 `json:"startTimestamp"`

	// Open timestamp in unix milliseconds.
	OpenTimestamp int64 `json:"openTimestamp"`

	// Close timestamp in unix milliseconds.
	CloseTimestamp int64 `json:"closeTimestamp"`
}

func (t *Ticker24hData) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		market         = util.GetOrEmpty[string]("market", j)
		open           = util.GetOrEmpty[string]("open", j)
		high           = util.GetOrEmpty[string]("high", j)
		low            = util.GetOrEmpty[string]("low", j)
		last           = util.GetOrEmpty[string]("last", j)
		volume         = util.GetOrEmpty[string]("volume", j)
		volumeQuote    = util.GetOrEmpty[string]("volumeQuote", j)
		bid            = util.GetOrEmpty[string]("bid", j)
		bidSize        = util.GetOrEmpty[string]("bidSize", j)
		ask            = util.GetOrEmpty[string]("ask", j)
		askSize        = util.GetOrEmpty[string]("askSize", j)
		timestamp      = util.GetOrEmpty[float64]("timestamp", j)
		startTimestamp = util.GetOrEmpty[float64]("startTimestamp", j)
		openTimestamp  = util.GetOrEmpty[float64]("openTimestamp", j)
		closeTimestamp = util.GetOrEmpty[float64]("closeTimestamp", j)
	)

	t.Market = market
	t.Open = util.IfOrElse(len(open) > 0, func() float64 { return util.MustFloat64(open) }, 0)
	t.High = util.IfOrElse(len(high) > 0, func() float64 { return util.MustFloat64(high) }, 0)
	t.Low = util.IfOrElse(len(low) > 0, func() float64 { return util.MustFloat64(low) }, 0)
	t.Last = util.IfOrElse(len(last) > 0, func() float64 { return util.MustFloat64(last) }, 0)
	t.Volume = util.IfOrElse(len(volume) > 0, func() float64 { return util.MustFloat64(volume) }, 0)
	t.VolumeQuote = util.IfOrElse(len(volumeQuote) > 0, func() float64 { return util.MustFloat64(volumeQuote) }, 0)
	t.Bid = util.IfOrElse(len(bid) > 0, func() float64 { return util.MustFloat64(bid) }, 0)
	t.BidSize = util.IfOrElse(len(bidSize) > 0, func() float64 { return util.MustFloat64(bidSize) }, 0)
	t.Ask = util.IfOrElse(len(ask) > 0, func() float64 { return util.MustFloat64(ask) }, 0)
	t.AskSize = util.IfOrElse(len(askSize) > 0, func() float64 { return util.MustFloat64(askSize) }, 0)
	t.Timestamp = int64(timestamp)
	t.StartTimestamp = int64(startTimestamp)
	t.OpenTimestamp = int64(openTimestamp)
	t.CloseTimestamp = int64(closeTimestamp)

	return nil
}
