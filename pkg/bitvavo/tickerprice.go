package bitvavo

import (
	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/util"
)

type TickerPrice struct {
	// The market you requested the latest trade price for.
	Market string `json:"market"`

	// The latest trade price for 1 base currency in quote currency for market. For example, 34243 Euro.
	Price float64 `json:"price"`
}

func (t *TickerPrice) UnmarshalJSON(bytes []byte) error {
	var j map[string]string

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		market = j["market"]
		price  = j["price"]
	)

	t.Market = market
	t.Price = util.IfOrElse(len(price) > 0, func() float64 { return util.MustFloat64(price) }, 0)

	return nil
}
