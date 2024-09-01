package bitvavo

type Ticker struct {
	// The market which was requested in the subscription.
	Market string `json:"market"`

	// The price of the best (highest) bid offer available, only sent when either bestBid or bestBidSize has changed.
	BestBid string `json:"bestBid"`

	// The size of the best (highest) bid offer available, only sent when either bestBid or bestBidSize has changed.
	BestBidSize string `json:"bestBidSize"`

	// The price of the best (lowest) ask offer available, only sent when either bestAsk or bestAskSize has changed.
	BestAsk string `json:"bestAsk"`

	// The size of the best (lowest) ask offer available, only sent when either bestAsk or bestAskSize has changed.
	BestAskSize string `json:"bestAskSize"`

	// The last price for which a trade has occurred, only sent when lastPrice has changed.
	LastPrice string `json:"lastPrice"`
}
