package bitvavo

type Balance struct {
	// Short version of asset name.
	Symbol string `json:"symbol"`

	// Balance freely available.
	Available string `json:"available"`

	// Balance currently placed onHold for open orders.
	InOrder string `json:"inOrder"`
}
