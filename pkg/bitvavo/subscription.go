package bitvavo

import (
	"github.com/orsinium-labs/enum"
)

type Channel enum.Member[string]

var (
	channel          = enum.NewBuilder[string, Channel]()
	ChannelAccount   = channel.Add(Channel{"account"})
	ChannelBook      = channel.Add(Channel{"book"})
	ChannelCandles   = channel.Add(Channel{"candles"})
	ChannelTrades    = channel.Add(Channel{"trades"})
	ChannelTicker    = channel.Add(Channel{"ticker"})
	ChannelTicker24h = channel.Add(Channel{"ticker24h"})
	channels         = channel.Enum()
)

type Interval enum.Member[string]

var (
	interval    = enum.NewBuilder[string, Interval]()
	Interval1m  = interval.Add(Interval{"1m"})
	Interval5m  = interval.Add(Interval{"5m"})
	Interval15m = interval.Add(Interval{"15m"})
	Interval30m = interval.Add(Interval{"30m"})
	Interval1h  = interval.Add(Interval{"1h"})
	Interval2h  = interval.Add(Interval{"2h"})
	Interval4h  = interval.Add(Interval{"4h"})
	Interval6h  = interval.Add(Interval{"6h"})
	Interval8h  = interval.Add(Interval{"8h"})
	intervals   = interval.Enum()
)

type Subscription struct {
	Markets   []string
	Intervals []Interval
	Channel   Channel
}

func NewSubscription(channel Channel, markets []string, intervals ...Interval) Subscription {
	if channel == ChannelCandles && len(intervals) == 0 {
		panic("must provide at least one interval for candles channel")
	}
	return Subscription{
		Channel:   channel,
		Markets:   markets,
		Intervals: intervals,
	}
}
