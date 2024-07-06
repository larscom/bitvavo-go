package bitvavo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/larscom/bitvavo-go/internal/util"
	"github.com/orsinium-labs/enum"
)

type OrderParams struct {
	// Return the limit most recent orders only.
	// Default: 500
	Limit uint64 `json:"limit"`

	// Return orders after start time.
	Start time.Time `json:"start"`

	// Return orders before end time.
	End time.Time `json:"end"`

	// Filter used to limit the returned results.
	// All orders after this order ID are returned (i.e. showing those later in time).
	OrderIdFrom string `json:"orderIdFrom"`

	// Filter used to limit the returned results.
	// All orders up to this order ID are returned (i.e. showing those earlier in time).
	OrderIdTo string `json:"orderIdTo"`
}

func (o *OrderParams) Params() url.Values {
	params := make(url.Values)
	if o.Limit > 0 {
		params.Add("limit", fmt.Sprint(o.Limit))
	}
	if !o.Start.IsZero() {
		params.Add("start", fmt.Sprint(o.Start.UnixMilli()))
	}
	if !o.End.IsZero() {
		params.Add("end", fmt.Sprint(o.End.UnixMilli()))
	}
	if o.OrderIdFrom != "" {
		params.Add("orderIdFrom", o.OrderIdFrom)
	}
	if o.OrderIdTo != "" {
		params.Add("orderIdTo", o.OrderIdTo)
	}
	return params
}

type OrderStatus enum.Member[string]

var (
	orderStatus                   = enum.NewBuilder[string, OrderStatus]()
	ORDER_STATUS_NEW              = orderStatus.Add(OrderStatus{"new"})
	ORDER_STATUS_AWAITING_TRIGGER = orderStatus.Add(OrderStatus{"awaitingTrigger"})
	ORDER_STATUS_CANCELED         = orderStatus.Add(OrderStatus{"canceled"})
	ORDER_STATUS_CANCELED_AUCTION = orderStatus.Add(OrderStatus{"canceledAuction"})
	ORDER_STATUS_CANCELED_STP     = orderStatus.Add(OrderStatus{"canceledSelfTradePrevention"})
	ORDER_STATUS_CANCELED_IOC     = orderStatus.Add(OrderStatus{"canceledIOC"})
	ORDER_STATUS_CANCELED_FOK     = orderStatus.Add(OrderStatus{"canceledFOK"})
	ORDER_STATUS_CANCELED_MP      = orderStatus.Add(OrderStatus{"canceledMarketProtection"})
	ORDER_STATUS_CANCELED_PO      = orderStatus.Add(OrderStatus{"canceledPostOnly"})
	ORDER_STATUS_FILLED           = orderStatus.Add(OrderStatus{"filled"})
	ORDER_STATUS_PARTIALLY_FILLED = orderStatus.Add(OrderStatus{"partiallyFilled"})
	ORDER_STATUS_EXPIRED          = orderStatus.Add(OrderStatus{"expired"})
	ORDER_STATUS_REJECTED         = orderStatus.Add(OrderStatus{"rejected"})
	orderStatuses                 = orderStatus.Enum()
)

type OrderType enum.Member[string]

var (
	orderType                    = enum.NewBuilder[string, OrderType]()
	ORDER_TYPE_MARKET            = orderType.Add(OrderType{"market"})
	ORDER_TYPE_LIMIT             = orderType.Add(OrderType{"limit"})
	ORDER_TYPE_STOP_LOSS         = orderType.Add(OrderType{"stopLoss"})
	ORDER_TYPE_STOP_LOSS_LIMIT   = orderType.Add(OrderType{"stopLossLimit"})
	ORDER_TYPE_TAKE_PROFIT       = orderType.Add(OrderType{"takeProfit"})
	ORDER_TYPE_TAKE_PROFIT_LIMIT = orderType.Add(OrderType{"takeProfitLimit"})
	orderTypes                   = orderType.Enum()
)

type OrderTriggerType enum.Member[string]

var (
	orderTriggerType           = enum.NewBuilder[string, OrderTriggerType]()
	ORDER_TRIGGER_TYPE_DEFAULT = ORDER_TRIGGER_TYPE_PRICE
	ORDER_TRIGGER_TYPE_PRICE   = orderTriggerType.Add(OrderTriggerType{"price"})
	orderTriggerTypes          = orderTriggerType.Enum()
)

type OrderTriggerRef enum.Member[string]

var (
	orderTriggerRef              = enum.NewBuilder[string, OrderTriggerRef]()
	ORDER_TRIGGER_REF_LAST_TRADE = orderTriggerRef.Add(OrderTriggerRef{"lastTrade"})
	ORDER_TRIGGER_REF_BEST_BID   = orderTriggerRef.Add(OrderTriggerRef{"bestBid"})
	ORDER_TRIGGER_REF_BEST_ASK   = orderTriggerRef.Add(OrderTriggerRef{"bestAsk"})
	ORDER_TRIGGER_REF_MID_PRICE  = orderTriggerRef.Add(OrderTriggerRef{"midPrice"})
	orderTriggerRefs             = orderTriggerRef.Enum()
)

type TimeInForce enum.Member[string]

var (
	timeInForce           = enum.NewBuilder[string, TimeInForce]()
	TIME_IN_FORCE_DEFAULT = TIME_IN_FORCE_GTC
	TIME_IN_FORCE_GTC     = timeInForce.Add(TimeInForce{"GTC"})
	TIME_IN_FORCE_IOC     = timeInForce.Add(TimeInForce{"IOC"})
	TIME_IN_FORCE_FOK     = timeInForce.Add(TimeInForce{"FOK"})
	timeInForces          = timeInForce.Enum()
)

type SelfTradePrevention enum.Member[string]

var (
	selfTradePrevention           = enum.NewBuilder[string, SelfTradePrevention]()
	SELF_TRADE_PREVENTION_DEFAULT = SELF_TRADE_PREVENTION_DAC
	SELF_TRADE_PREVENTION_DAC     = selfTradePrevention.Add(SelfTradePrevention{"decrementAndCancel"})
	SELF_TRADE_PREVENTION_CO      = selfTradePrevention.Add(SelfTradePrevention{"cancelOldest"})
	SELF_TRADE_PREVENTION_CN      = selfTradePrevention.Add(SelfTradePrevention{"cancelNewest"})
	SELF_TRADE_PREVENTION_CB      = selfTradePrevention.Add(SelfTradePrevention{"cancelBoth"})
	selfTradePreventions          = selfTradePrevention.Enum()
)

type Order struct {
	// The order id of the returned order.
	OrderId string `json:"orderId"`

	// The market in which the order was placed.
	Market string `json:"market"`

	// Is a timestamp in milliseconds since 1 Jan 1970.
	Created int64 `json:"created"`

	// Is a timestamp in milliseconds since 1 Jan 1970.
	Updated int64 `json:"updated"`

	// The current status of the order.
	Status OrderStatus `json:"status"`

	// Side
	Side Side `json:"side"`

	// OrderType
	OrderType OrderType `json:"orderType"`

	// Original amount.
	Amount float64 `json:"amount"`

	// Amount remaining (lower than 'amount' after fills).
	AmountRemaining float64 `json:"amountRemaining"`

	// The price of the order.
	Price float64 `json:"price"`

	// Amount of 'onHoldCurrency' that is reserved for this order. This is released when orders are canceled.
	OnHold float64 `json:"onHold"`

	// The currency placed on hold is the quote currency for sell orders and base currency for buy orders.
	OnHoldCurrency string `json:"onHoldCurrency"`

	// Only for stop orders: The current price used in the trigger. This is based on the triggerAmount and triggerType.
	TriggerPrice float64 `json:"triggerPrice"`

	// Only for stop orders: The value used for the triggerType to determine the triggerPrice.
	TriggerAmount float64 `json:"triggerAmount"`

	// Only for stop orders.
	TriggerType OrderTriggerType `json:"triggerType"`

	// Only for stop orders: The reference price used for stop orders.
	TriggerReference OrderTriggerRef `json:"triggerReference"`

	// Only for limit orders: Determines how long orders remain active.
	// Possible values: Good-Til-Canceled (GTC), Immediate-Or-Cancel (IOC), Fill-Or-Kill (FOK).
	// GTC orders will remain on the order book until they are filled or canceled.
	// IOC orders will fill against existing orders, but will cancel any remaining amount after that.
	// FOK orders will fill against existing orders in its entirety, or will be canceled (if the entire order cannot be filled).
	//
	TimeInForce TimeInForce `json:"timeInForce"`

	// Default: false
	PostOnly bool `json:"postOnly"`

	// Self trading is not allowed on Bitvavo. Multiple options are available to prevent this from happening.
	// The default ‘decrementAndCancel’ decrements both orders by the amount that would have been filled, which in turn cancels the smallest of the two orders.
	// ‘cancelOldest’ will cancel the entire older order and places the new order.
	// ‘cancelNewest’ will cancel the order that is submitted.
	// ‘cancelBoth’ will cancel both the current and the old order.
	//
	// Default: "decrementAndCancel"
	SelfTradePrevention SelfTradePrevention `json:"selfTradePrevention"`

	// Whether this order is visible on the order book.
	Visible bool `json:"visible"`

	// The fills for this order
	Fills []Fill `json:"fills"`

	// How much of this order is filled
	FilledAmount float64 `json:"filledAmount"`

	// How much of this order is filled in quote currency
	FilledAmountQuote float64 `json:"filledAmountQuote"`

	// The currency in which the fee is paid (e.g: EUR)
	FeeCurrency string `json:"feeCurrency"`

	// How much fee is paid
	FeePaid float64 `json:"feePaid"`
}

func (o *Order) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		orderId             = util.GetOrEmpty[string]("orderId", j)
		market              = util.GetOrEmpty[string]("market", j)
		created             = util.GetOrEmpty[float64]("created", j)
		updated             = util.GetOrEmpty[float64]("updated", j)
		status              = util.GetOrEmpty[string]("status", j)
		side                = util.GetOrEmpty[string]("side", j)
		orderType           = util.GetOrEmpty[string]("orderType", j)
		amount              = util.GetOrEmpty[string]("amount", j)
		amountRemaining     = util.GetOrEmpty[string]("amountRemaining", j)
		price               = util.GetOrEmpty[string]("price", j)
		onHold              = util.GetOrEmpty[string]("onHold", j)
		onHoldCurrency      = util.GetOrEmpty[string]("onHoldCurrency", j)
		timeInForce         = util.GetOrEmpty[string]("timeInForce", j)
		postOnly            = util.GetOrEmpty[bool]("postOnly", j)
		selfTradePrevention = util.GetOrEmpty[string]("selfTradePrevention", j)
		visible             = util.GetOrEmpty[bool]("visible", j)

		// only for stop orders
		triggerPrice     = util.GetOrEmpty[string]("triggerPrice", j)
		triggerAmount    = util.GetOrEmpty[string]("triggerAmount", j)
		triggerType      = util.GetOrEmpty[string]("triggerType", j)
		triggerReference = util.GetOrEmpty[string]("triggerReference", j)

		fillsAny          = util.GetOrEmpty[[]any]("fills", j)
		filledAmount      = util.GetOrEmpty[string]("filledAmount", j)
		filledAmountQuote = util.GetOrEmpty[string]("filledAmountQuote", j)
		feeCurrency       = util.GetOrEmpty[string]("feeCurrency", j)
		feePaid           = util.GetOrEmpty[string]("feePaid", j)
	)

	if len(fillsAny) > 0 {
		fillsBytes, err := json.Marshal(fillsAny)
		if err != nil {
			return err
		}
		fills := make([]Fill, len(fillsAny))
		if err := json.Unmarshal(fillsBytes, &fills); err != nil {
			return err
		}
		o.Fills = fills
	}

	o.OrderId = orderId
	o.Market = market
	o.Created = int64(created)
	o.Updated = int64(updated)
	o.Status = *orderStatuses.Parse(status)
	o.Side = *sides.Parse(side)
	o.OrderType = *orderTypes.Parse(orderType)
	o.Amount = util.IfOrElse(len(amount) > 0, func() float64 { return util.MustFloat64(amount) }, 0)
	o.AmountRemaining = util.IfOrElse(len(amountRemaining) > 0, func() float64 { return util.MustFloat64(amountRemaining) }, 0)
	o.Price = util.IfOrElse(len(price) > 0, func() float64 { return util.MustFloat64(price) }, 0)
	o.OnHold = util.IfOrElse(len(onHold) > 0, func() float64 { return util.MustFloat64(onHold) }, 0)
	o.OnHoldCurrency = onHoldCurrency
	o.TriggerPrice = util.IfOrElse(len(triggerPrice) > 0, func() float64 { return util.MustFloat64(triggerPrice) }, 0)
	o.TriggerAmount = util.IfOrElse(len(triggerAmount) > 0, func() float64 { return util.MustFloat64(triggerAmount) }, 0)
	if len(triggerType) > 0 {
		o.TriggerType = *orderTriggerTypes.Parse(triggerType)
	}
	if len(triggerReference) > 0 {
		o.TriggerReference = *orderTriggerRefs.Parse(triggerReference)
	}
	if len(timeInForce) > 0 {
		o.TimeInForce = *timeInForces.Parse(timeInForce)
	}
	o.PostOnly = postOnly
	o.SelfTradePrevention = *selfTradePreventions.Parse(selfTradePrevention)
	o.Visible = visible
	o.FilledAmount = util.IfOrElse(len(filledAmount) > 0, func() float64 { return util.MustFloat64(filledAmount) }, 0)
	o.FilledAmountQuote = util.IfOrElse(len(filledAmountQuote) > 0, func() float64 { return util.MustFloat64(filledAmountQuote) }, 0)
	o.FeeCurrency = feeCurrency
	o.FeePaid = util.IfOrElse(len(feePaid) > 0, func() float64 { return util.MustFloat64(feePaid) }, 0)

	return nil
}

type OrderNew struct {
	// The market in which the order should be placed (e.g: ETH-EUR)
	Market string `json:"market"`

	// When placing a buy order the base currency will be bought for the quote currency. When placing a sell order the base currency will be sold for the quote currency.
	Side Side `json:"side"`

	// For limit orders, amount and price are required. For market orders either amount or amountQuote is required.
	OrderType OrderType `json:"orderType"`

	// Specifies the amount of the base asset that will be bought/sold.
	Amount float64 `json:"amount,omitempty"`

	// Only for limit orders: Specifies the amount in quote currency that is paid/received for each unit of base currency.
	Price float64 `json:"price,omitempty"`

	// Only for market orders: If amountQuote is specified, [amountQuote] of the quote currency will be bought/sold for the best price available.
	AmountQuote float64 `json:"amountQuote,omitempty"`

	// Only for stop orders: Specifies the amount that is used with the triggerType.
	// Combine this parameter with triggerType and triggerReference to create the desired trigger.
	TriggerAmount float64 `json:"triggerAmount,omitempty"`

	// Only for stop orders: Only allows price for now. A triggerAmount of 4000 and a triggerType of price will generate a triggerPrice of 4000.
	// Combine this parameter with triggerAmount and triggerReference to create the desired trigger.
	TriggerType OrderTriggerType `json:"triggerType,omitempty"`

	// Only for stop orders: Use this to determine which parameter will trigger the order.
	// Combine this parameter with triggerAmount and triggerType to create the desired trigger.
	TriggerReference OrderTriggerRef `json:"triggerReference,omitempty"`

	// Only for limit orders: Determines how long orders remain active.
	// Possible values: Good-Til-Canceled (GTC), Immediate-Or-Cancel (IOC), Fill-Or-Kill (FOK).
	// GTC orders will remain on the order book until they are filled or canceled.
	// IOC orders will fill against existing orders, but will cancel any remaining amount after that.
	// FOK orders will fill against existing orders in its entirety, or will be canceled (if the entire order cannot be filled).
	//
	// Default: "GTC"
	TimeInForce TimeInForce `json:"timeInForce,omitempty"`

	// Self trading is not allowed on Bitvavo. Multiple options are available to prevent this from happening.
	// The default ‘decrementAndCancel’ decrements both orders by the amount that would have been filled, which in turn cancels the smallest of the two orders.
	// ‘cancelOldest’ will cancel the entire older order and places the new order.
	// ‘cancelNewest’ will cancel the order that is submitted.
	// ‘cancelBoth’ will cancel both the current and the old order.
	//
	// Default: "decrementAndCancel"
	SelfTradePrevention SelfTradePrevention `json:"selfTradePrevention,omitempty"`

	// Only for limit orders: When postOnly is set to true, the order will not fill against existing orders.
	// This is useful if you want to ensure you pay the maker fee. If the order would fill against existing orders, the entire order will be canceled.
	//
	// Default: false
	PostOnly bool `json:"postOnly,omitempty"`

	// Only for market orders: In order to protect clients from filling market orders with undesirable prices,
	// the remainder of market orders will be canceled once the next fill price is 10% worse than the best fill price (best bid/ask at first match).
	// If you wish to disable this protection, set this value to ‘true’.
	//
	// Default: false
	DisableMarketProtection bool `json:"disableMarketProtection,omitempty"`

	// If this is set to 'true', all order information is returned.
	// Set this to 'false' when only an acknowledgement of success or failure is required, this is faster.
	//
	// Default: true
	ResponseRequired bool `json:"responseRequired,omitempty"`
}

func (o OrderNew) MarshalJSON() ([]byte, error) {
	type O OrderNew

	target := &struct {
		O
		Side                string `json:"side"`
		OrderType           string `json:"orderType"`
		TriggerType         string `json:"triggerType,omitempty"`
		TriggerReference    string `json:"triggerReference,omitempty"`
		TimeInForce         string `json:"timeInForce,omitempty"`
		SelfTradePrevention string `json:"selfTradePrevention,omitempty"`
	}{
		O:                   (O)(o),
		Side:                o.Side.Value,
		OrderType:           o.OrderType.Value,
		TriggerType:         o.TriggerType.Value,
		TriggerReference:    o.TriggerReference.Value,
		TimeInForce:         o.TimeInForce.Value,
		SelfTradePrevention: o.SelfTradePrevention.Value,
	}

	return json.Marshal(target)
}

type OrderUpdate struct {
	// The market for which an order should be updated
	Market string `json:"market"`

	// The id of the order which should be updated
	OrderId string `json:"orderId"`

	// Updates amount to this value (and also changes amountRemaining accordingly).
	Amount float64 `json:"amount,omitempty"`

	// Only for market orders: If amountQuote is specified, [amountQuote] of the quote currency will be bought/sold for the best price available.
	AmountQuote float64 `json:"amountQuote,omitempty"`

	// Updates amountRemaining to this value (and also changes amount accordingly).
	AmountRemaining float64 `json:"amountRemaining,omitempty"`

	// Specifies the amount in quote currency that is paid/received for each unit of base currency.
	Price float64 `json:"price,omitempty"`

	// Only for stop orders: Specifies the amount that is used with the triggerType.
	// Combine this parameter with triggerType and triggerReference to create the desired trigger.
	TriggerAmount float64 `json:"triggerAmount,omitempty"`

	// Only for limit orders: Determines how long orders remain active.
	// Possible values: Good-Til-Canceled (GTC), Immediate-Or-Cancel (IOC), Fill-Or-Kill (FOK).
	// GTC orders will remain on the order book until they are filled or canceled.
	// IOC orders will fill against existing orders, but will cancel any remaining amount after that.
	// FOK orders will fill against existing orders in its entirety, or will be canceled (if the entire order cannot be filled).
	//
	// Default: "GTC"
	TimeInForce TimeInForce `json:"timeInForce,omitempty"`

	// Self trading is not allowed on Bitvavo. Multiple options are available to prevent this from happening.
	// The default ‘decrementAndCancel’ decrements both orders by the amount that would have been filled, which in turn cancels the smallest of the two orders.
	// ‘cancelOldest’ will cancel the entire older order and places the new order.
	// ‘cancelNewest’ will cancel the order that is submitted.
	// ‘cancelBoth’ will cancel both the current and the old order.
	//
	// Default: "decrementAndCancel"
	SelfTradePrevention SelfTradePrevention `json:"selfTradePrevention,omitempty"`

	// Only for limit orders: When postOnly is set to true, the order will not fill against existing orders.
	// This is useful if you want to ensure you pay the maker fee. If the order would fill against existing orders, the entire order will be canceled.
	//
	// Default: false
	PostOnly bool `json:"postOnly,omitempty"`

	// If this is set to 'true', all order information is returned.
	// Set this to 'false' when only an acknowledgement of success or failure is required, this is faster.
	//
	// Default: true
	ResponseRequired bool `json:"responseRequired,omitempty"`
}

func (o OrderUpdate) MarshalJSON() ([]byte, error) {
	type O OrderUpdate

	target := &struct {
		O
		TimeInForce         string `json:"timeInForce,omitempty"`
		SelfTradePrevention string `json:"selfTradePrevention,omitempty"`
	}{
		O:                   (O)(o),
		TimeInForce:         o.TimeInForce.Value,
		SelfTradePrevention: o.SelfTradePrevention.Value,
	}

	return json.Marshal(target)
}
