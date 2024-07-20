package bitvavo

import "github.com/orsinium-labs/enum"

type Side enum.Member[string]

var (
	side     = enum.NewBuilder[string, Side]()
	SideBuy  = side.Add(Side{"buy"})
	SideSell = side.Add(Side{"sell"})
	sides    = side.Enum()
)
