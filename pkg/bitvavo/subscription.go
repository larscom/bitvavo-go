package bitvavo

import (
	"github.com/orsinium-labs/enum"
)

type Channel enum.Member[string]

var (
	channel           = enum.NewBuilder[string, Channel]()
	CHANNEL_ACCOUNT   = channel.Add(Channel{"account"})
	CHANNEL_BOOK      = channel.Add(Channel{"book"})
	CHANNEL_CANDLES   = channel.Add(Channel{"candles"})
	CHANNEL_TRADES    = channel.Add(Channel{"trades"})
	CHANNEL_TICKER    = channel.Add(Channel{"ticker"})
	CHANNEL_TICKER24H = channel.Add(Channel{"ticker24h"})
	channels          = channel.Enum()
)

type Interval enum.Member[string]

var (
	interval     = enum.NewBuilder[string, Interval]()
	INTERVAL_1M  = interval.Add(Interval{"1m"})
	INTERVAL_5M  = interval.Add(Interval{"5m"})
	INTERVAL_15M = interval.Add(Interval{"15m"})
	INTERVAL_30M = interval.Add(Interval{"30m"})
	INTERVAL_1H  = interval.Add(Interval{"1h"})
	INTERVAL_2H  = interval.Add(Interval{"2h"})
	INTERVAL_4H  = interval.Add(Interval{"4h"})
	INTERVAL_6H  = interval.Add(Interval{"6h"})
	INTERVAL_8H  = interval.Add(Interval{"8h"})
	intervals    = interval.Enum()
)

type Subscription struct {
	Markets   []string
	Intervals []Interval
	Channel   Channel
}

func NewSubscription(channel Channel, markets []string, intervals ...Interval) Subscription {
	if channel == CHANNEL_CANDLES && len(intervals) == 0 {
		panic("must provide at least one interval for candles channel")
	}
	return Subscription{
		Channel:   channel,
		Markets:   markets,
		Intervals: intervals,
	}
}
