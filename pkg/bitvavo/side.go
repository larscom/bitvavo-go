package bitvavo

import "github.com/orsinium-labs/enum"

type Side enum.Member[string]

var (
	side      = enum.NewBuilder[string, Side]()
	SIDE_BUY  = side.Add(Side{"buy"})
	SIDE_SELL = side.Add(Side{"sell"})
	sides     = side.Enum()
)
