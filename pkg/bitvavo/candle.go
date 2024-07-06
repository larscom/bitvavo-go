package bitvavo

import (
	"fmt"
	"net/url"
	"time"

	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/util"
)

var ErrExpectedCandleLenght = func(exp, act int) error { return fmt.Errorf("expected length '%d' for candle, but was: %d", exp, act) }

type CandleParams struct {
	// Return the limit most recent candlesticks only.
	// Default: 1440
	Limit uint64 `json:"limit"`

	// Return limit candlesticks for trades made after start.
	Start time.Time `json:"start"`

	// Return limit candlesticks for trades made before end.
	End time.Time `json:"end"`
}

func (c *CandleParams) Params() url.Values {
	params := make(url.Values)
	if c.Limit > 0 {
		params.Add("limit", fmt.Sprint(c.Limit))
	}
	if !c.Start.IsZero() {
		params.Add("start", fmt.Sprint(c.Start.UnixMilli()))
	}
	if !c.End.IsZero() {
		params.Add("end", fmt.Sprint(c.End.UnixMilli()))
	}
	return params
}

type CandleOnly struct {
	Timestamp int64   `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

func (c *CandleOnly) UnmarshalJSON(bytes []byte) error {
	var candle []any

	if err := json.Unmarshal(bytes, &candle); err != nil {
		return err
	}

	if len(candle) != 6 {
		return ErrExpectedCandleLenght(6, len(candle))
	}

	c.Timestamp = int64(candle[0].(float64))
	c.Open = util.MustFloat64(candle[1].(string))
	c.High = util.MustFloat64(candle[2].(string))
	c.Low = util.MustFloat64(candle[3].(string))
	c.Close = util.MustFloat64(candle[4].(string))
	c.Volume = util.MustFloat64(candle[5].(string))

	return nil
}

type Candle struct {
	Interval  Interval `json:"interval"`
	Market    string   `json:"market"`
	Timestamp int64    `json:"timestamp"`
	Open      float64  `json:"open"`
	High      float64  `json:"high"`
	Low       float64  `json:"low"`
	Close     float64  `json:"close"`
	Volume    float64  `json:"volume"`
}

func (c *Candle) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		market   = util.GetOrEmpty[string]("market", j)
		interval = util.GetOrEmpty[string]("interval", j)
		candles  = util.GetOrEmpty[[]any]("candle", j)
	)

	if len(candles) != 1 {
		return ErrExpectedCandleLenght(1, len(candles))
	}

	candle := candles[0].([]any)
	if len(candle) != 6 {
		return ErrExpectedCandleLenght(6, len(candle))
	}

	c.Market = market
	c.Interval = *intervals.Parse(interval)
	c.Timestamp = int64(candle[0].(float64))
	c.Open = util.MustFloat64(candle[1].(string))
	c.High = util.MustFloat64(candle[2].(string))
	c.Low = util.MustFloat64(candle[3].(string))
	c.Close = util.MustFloat64(candle[4].(string))
	c.Volume = util.MustFloat64(candle[5].(string))

	return nil
}
