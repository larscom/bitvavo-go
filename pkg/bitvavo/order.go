package bitvavo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/larscom/bitvavo-go/v2/internal/util"
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
	orderStatus                = enum.NewBuilder[string, OrderStatus]()
	OrderStatusNew             = orderStatus.Add(OrderStatus{"new"})
	OrderStatusAwaitingTrigger = orderStatus.Add(OrderStatus{"awaitingTrigger"})
	OrderStatusCanceled        = orderStatus.Add(OrderStatus{"canceled"})
	OrderStatusCanceledAuction = orderStatus.Add(OrderStatus{"canceledAuction"})
	OrderStatusCanceledStp     = orderStatus.Add(OrderStatus{"canceledSelfTradePrevention"})
	OrderStatusCanceledIoc     = orderStatus.Add(OrderStatus{"canceledIOC"})
	OrderStatusCanceledFok     = orderStatus.Add(OrderStatus{"canceledFOK"})
	OrderStatusCanceledMp      = orderStatus.Add(OrderStatus{"canceledMarketProtection"})
	OrderStatusCanceledPo      = orderStatus.Add(OrderStatus{"canceledPostOnly"})
	OrderStatusFilled          = orderStatus.Add(OrderStatus{"filled"})
	OrderStatusPartiallyFilled = orderStatus.Add(OrderStatus{"partiallyFilled"})
	OrderStatusExpired         = orderStatus.Add(OrderStatus{"expired"})
	OrderStatusRejected        = orderStatus.Add(OrderStatus{"rejected"})
	orderStatuses              = orderStatus.Enum()
)

type OrderType enum.Member[string]

var (
	orderType                = enum.NewBuilder[string, OrderType]()
	OrderTypeMarket          = orderType.Add(OrderType{"market"})
	OrderTypeLimit           = orderType.Add(OrderType{"limit"})
	OrderTypeStopLoss        = orderType.Add(OrderType{"stopLoss"})
	OrderTypeStopLossLimit   = orderType.Add(OrderType{"stopLossLimit"})
	OrderTypeTakeProfit      = orderType.Add(OrderType{"takeProfit"})
	OrderTypeTakeProfitLimit = orderType.Add(OrderType{"takeProfitLimit"})
	orderTypes               = orderType.Enum()
)

type OrderTriggerType enum.Member[string]

var (
	orderTriggerType        = enum.NewBuilder[string, OrderTriggerType]()
	OrderTriggerTypeDefault = OrderTriggerTypePrice
	OrderTriggerTypePrice   = orderTriggerType.Add(OrderTriggerType{"price"})
	orderTriggerTypes       = orderTriggerType.Enum()
)

type OrderTriggerRef enum.Member[string]

var (
	orderTriggerRef          = enum.NewBuilder[string, OrderTriggerRef]()
	OrderTriggerRefLastTrade = orderTriggerRef.Add(OrderTriggerRef{"lastTrade"})
	OrderTriggerRefBestBid   = orderTriggerRef.Add(OrderTriggerRef{"bestBid"})
	OrderTriggerRefBestAsk   = orderTriggerRef.Add(OrderTriggerRef{"bestAsk"})
	OrderTriggerRefMidPrice  = orderTriggerRef.Add(OrderTriggerRef{"midPrice"})
	orderTriggerRefs         = orderTriggerRef.Enum()
)

type TimeInForce enum.Member[string]

var (
	timeInForce        = enum.NewBuilder[string, TimeInForce]()
	TimeInForceDefault = TimeInForceGtc
	TimeInForceGtc     = timeInForce.Add(TimeInForce{"GTC"})
	TimeInForceIoc     = timeInForce.Add(TimeInForce{"IOC"})
	TimeInForceFok     = timeInForce.Add(TimeInForce{"FOK"})
	timeInForces       = timeInForce.Enum()
)

type SelfTradePrevention enum.Member[string]

var (
	selfTradePrevention        = enum.NewBuilder[string, SelfTradePrevention]()
	SelfTradePreventionDefault = SelfTradePreventionDac
	SelfTradePreventionDac     = selfTradePrevention.Add(SelfTradePrevention{"decrementAndCancel"})
	SelfTradePreventionCo      = selfTradePrevention.Add(SelfTradePrevention{"cancelOldest"})
	SelfTradePreventionCn      = selfTradePrevention.Add(SelfTradePrevention{"cancelNewest"})
	SelfTradePreventionCb      = selfTradePrevention.Add(SelfTradePrevention{"cancelBoth"})
	selfTradePreventions       = selfTradePrevention.Enum()
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
	Amount string `json:"amount"`

	// Amount remaining (lower than 'amount' after fills).
	AmountRemaining string `json:"amountRemaining"`

	// The price of the order.
	Price string `json:"price"`

	// Amount of 'onHoldCurrency' that is reserved for this order. This is released when orders are canceled.
	OnHold string `json:"onHold"`

	// The currency placed on hold is the quote currency for sell orders and base currency for buy orders.
	OnHoldCurrency string `json:"onHoldCurrency"`

	// Only for stop orders: The current price used in the trigger. This is based on the triggerAmount and triggerType.
	TriggerPrice string `json:"triggerPrice"`

	// Only for stop orders: The value used for the triggerType to determine the triggerPrice.
	TriggerAmount string `json:"triggerAmount"`

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
	FilledAmount string `json:"filledAmount"`

	// How much of this order is filled in quote currency
	FilledAmountQuote string `json:"filledAmountQuote"`

	// The currency in which the fee is paid (e.g: EUR)
	FeeCurrency string `json:"feeCurrency"`

	// How much fee is paid
	FeePaid string `json:"feePaid"`
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
	o.Amount = amount
	o.AmountRemaining = amountRemaining
	o.Price = price
	o.OnHold = onHold
	o.OnHoldCurrency = onHoldCurrency
	o.TriggerPrice = triggerPrice
	o.TriggerAmount = triggerAmount
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
	o.FilledAmount = filledAmount
	o.FilledAmountQuote = filledAmountQuote
	o.FeeCurrency = feeCurrency
	o.FeePaid = feePaid

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
	Amount string `json:"amount,omitempty"`

	// Only for limit orders: Specifies the amount in quote currency that is paid/received for each unit of base currency.
	Price string `json:"price,omitempty"`

	// Only for market orders: If amountQuote is specified, [amountQuote] of the quote currency will be bought/sold for the best price available.
	AmountQuote string `json:"amountQuote,omitempty"`

	// Only for stop orders: Specifies the amount that is used with the triggerType.
	// Combine this parameter with triggerType and triggerReference to create the desired trigger.
	TriggerAmount string `json:"triggerAmount,omitempty"`

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
	Amount string `json:"amount,omitempty"`

	// Only for market orders: If amountQuote is specified, [amountQuote] of the quote currency will be bought/sold for the best price available.
	AmountQuote string `json:"amountQuote,omitempty"`

	// Updates amountRemaining to this value (and also changes amount accordingly).
	AmountRemaining string `json:"amountRemaining,omitempty"`

	// Specifies the amount in quote currency that is paid/received for each unit of base currency.
	Price string `json:"price,omitempty"`

	// Only for stop orders: Specifies the amount that is used with the triggerType.
	// Combine this parameter with triggerType and triggerReference to create the desired trigger.
	TriggerAmount string `json:"triggerAmount,omitempty"`

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
