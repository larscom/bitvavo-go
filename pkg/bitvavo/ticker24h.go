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
	Open string `json:"open"`

	// The highest price for which a trade occurred in the 24-hour period.
	High string `json:"high"`

	// The lowest price for which a trade occurred in the 24-hour period.
	Low string `json:"low"`

	// The last price for which a trade occurred in the 24-hour period.
	Last string `json:"last"`

	// The total volume of the 24-hour period in base currency.
	Volume string `json:"volume"`

	// The total volume of the 24-hour period in quote currency.
	VolumeQuote string `json:"volumeQuote"`

	// The best (highest) bid offer at the current moment.
	Bid string `json:"bid"`

	// The size of the best (highest) bid offer.
	BidSize string `json:"bidSize"`

	// The best (lowest) ask offer at the current moment.
	Ask string `json:"ask"`

	// The size of the best (lowest) ask offer.
	AskSize string `json:"askSize"`

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
	t.Open = open
	t.High = high
	t.Low = low
	t.Last = last
	t.Volume = volume
	t.VolumeQuote = volumeQuote
	t.Bid = bid
	t.BidSize = bidSize
	t.Ask = ask
	t.AskSize = askSize
	t.Timestamp = int64(timestamp)
	t.StartTimestamp = int64(startTimestamp)
	t.OpenTimestamp = int64(openTimestamp)
	t.CloseTimestamp = int64(closeTimestamp)

	return nil
}
