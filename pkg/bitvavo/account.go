package bitvavo

type Account struct {
	Fees Fee `json:"fees"`
}

type Fee struct {
	// Fee for trades that take liquidity from the order book.
	Taker string `json:"taker"`

	// Fee for trades that add liquidity to the order book.
	Maker string `json:"maker"`

	// Your trading volume in the last 30 days measured in EUR.
	Volume string `json:"volume"`
}
