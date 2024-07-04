package bitvavo

import (
	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/util"
)

type Balance struct {
	// Short version of asset name.
	Symbol string `json:"symbol"`

	// Balance freely available.
	Available float64 `json:"available"`

	// Balance currently placed onHold for open orders.
	InOrder float64 `json:"inOrder"`
}

func (b *Balance) UnmarshalJSON(bytes []byte) error {
	var j map[string]string

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		symbol    = j["symbol"]
		available = j["available"]
		inOrder   = j["inOrder"]
	)

	b.Symbol = symbol
	b.Available = util.IfOrElse(len(available) > 0, func() float64 { return util.MustFloat64(available) }, 0)
	b.InOrder = util.IfOrElse(len(inOrder) > 0, func() float64 { return util.MustFloat64(inOrder) }, 0)

	return nil
}
