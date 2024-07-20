package bitvavo

import (
	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/util"
	"github.com/orsinium-labs/enum"
)

type MarketStatus enum.Member[string]

var (
	marketStatus        = enum.NewBuilder[string, MarketStatus]()
	MarketStatusTrading = marketStatus.Add(MarketStatus{"trading"})
	MarketStatusHalted  = marketStatus.Add(MarketStatus{"halted"})
	MarketStatusAuction = marketStatus.Add(MarketStatus{"auction"})
	marketStatuses      = marketStatus.Enum()
)

type Market struct {
	// The market itself
	Market string `json:"market"`

	// The status of the market
	Status MarketStatus `json:"status"`

	// Base currency, found on the left side of the dash in market.
	Base string `json:"base"`

	// Quote currency, found on the right side of the dash in market.
	Quote string `json:"quote"`

	// Price precision determines how many significant digits are allowed. The rationale behind this is that for higher amounts, smaller price increments are less relevant.
	// Examples of valid prices for precision 5 are: 100010, 11313, 7500.10, 7500.20, 500.12, 0.0012345.
	// Examples of precision 6 are: 11313.1, 7500.11, 7500.25, 500.123, 0.00123456.
	PricePrecision int64 `json:"pricePrecision"`

	// The minimum amount in quote currency (amountQuote or amount * price) for valid orders.
	MinOrderInBaseAsset float64 `json:"minOrderInBaseAsset"`

	// The minimum amount in base currency for valid orders.
	MinOrderInQuoteAsset float64 `json:"minOrderInQuoteAsset"`

	// The maximum amount in quote currency (amountQuote or amount * price) for valid orders.
	MaxOrderInBaseAsset float64 `json:"maxOrderInBaseAsset"`

	// The maximum amount in base currency for valid orders.
	MaxOrderInQuoteAsset float64 `json:"maxOrderInQuoteAsset"`

	// Allowed order types for this market.
	OrderTypes []OrderType `json:"orderTypes"`
}

func (m *Market) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		market               = util.GetOrEmpty[string]("market", j)
		status               = util.GetOrEmpty[string]("status", j)
		base                 = util.GetOrEmpty[string]("base", j)
		quote                = util.GetOrEmpty[string]("quote", j)
		pricePrecision       = util.GetOrEmpty[float64]("pricePrecision", j)
		minOrderInBaseAsset  = util.GetOrEmpty[string]("minOrderInBaseAsset", j)
		minOrderInQuoteAsset = util.GetOrEmpty[string]("minOrderInQuoteAsset", j)
		maxOrderInBaseAsset  = util.GetOrEmpty[string]("maxOrderInBaseAsset", j)
		maxOrderInQuoteAsset = util.GetOrEmpty[string]("maxOrderInQuoteAsset", j)
		orderTypesAny        = util.GetOrEmpty[[]any]("orderTypes", j)
	)

	types := make([]OrderType, len(orderTypesAny))
	for i := 0; i < len(orderTypesAny); i++ {
		types[i] = *orderTypes.Parse(orderTypesAny[i].(string))
	}

	m.Market = market
	m.Status = *marketStatuses.Parse(status)
	m.Base = base
	m.Quote = quote
	m.PricePrecision = int64(pricePrecision)
	m.MinOrderInBaseAsset = util.IfOrElse(len(minOrderInBaseAsset) > 0, func() float64 { return util.MustFloat64(minOrderInBaseAsset) }, 0)
	m.MinOrderInQuoteAsset = util.IfOrElse(len(minOrderInQuoteAsset) > 0, func() float64 { return util.MustFloat64(minOrderInQuoteAsset) }, 0)
	m.MaxOrderInBaseAsset = util.IfOrElse(len(maxOrderInBaseAsset) > 0, func() float64 { return util.MustFloat64(maxOrderInBaseAsset) }, 0)
	m.MaxOrderInQuoteAsset = util.IfOrElse(len(maxOrderInQuoteAsset) > 0, func() float64 { return util.MustFloat64(maxOrderInQuoteAsset) }, 0)
	m.OrderTypes = types

	return nil
}
