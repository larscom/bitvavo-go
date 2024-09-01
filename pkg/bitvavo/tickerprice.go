package bitvavo

type TickerPrice struct {
	// The market you requested the latest trade price for.
	Market string `json:"market"`

	// The latest trade price for 1 base currency in quote currency for market. For example, 34243 Euro.
	Price string `json:"price"`
}
