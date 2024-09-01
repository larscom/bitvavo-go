package bitvavo

type TickerBook struct {
	// The market you requested the current best orders for.
	Market string `json:"market"`

	// The highest buy order in quote currency for market currently available on Bitvavo.
	Bid string `json:"bid"`

	// The amount of base currency for bid in the order.
	BidSize string `json:"bidSize"`

	// The lowest sell order in quote currency for market currently available on Bitvavo.
	Ask string `json:"ask"`

	// The amount of base currency for ask in the order.
	AskSize string `json:"askSize"`
}
