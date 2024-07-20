package bitvavo

import (
	"fmt"
	"github.com/larscom/bitvavo-go/internal/test"
	"testing"

	"github.com/goccy/go-json"
)

func TestMaxOrderNewMarshaller(t *testing.T) {
	order := OrderNew{
		Market:                  "ETH-EUR",
		Side:                    SideBuy,
		OrderType:               OrderTypeLimit,
		Amount:                  1.5,
		Price:                   2500.50,
		AmountQuote:             105.5,
		TriggerAmount:           10.2,
		TriggerType:             OrderTriggerTypeDefault,
		TriggerReference:        OrderTriggerRefBestAsk,
		TimeInForce:             TimeInForceDefault,
		SelfTradePrevention:     SelfTradePreventionDefault,
		PostOnly:                true,
		DisableMarketProtection: false,
		ResponseRequired:        true,
	}

	bytes, err := json.Marshal(order)
	if err != nil {
		t.Error(err)
	}

	expected := "{\"market\":\"ETH-EUR\",\"amount\":1.5,\"price\":2500.5,\"amountQuote\":105.5,\"triggerAmount\":10.2,\"postOnly\":true,\"responseRequired\":true,\"side\":\"buy\",\"orderType\":\"limit\",\"triggerType\":\"price\",\"triggerReference\":\"bestAsk\",\"timeInForce\":\"GTC\",\"selfTradePrevention\":\"decrementAndCancel\"}"
	actual := string(bytes)

	fmt.Printf("%q", actual)

	test.AssertEqual(t, expected, actual)
}

func TestMinOrderNewMarshaller(t *testing.T) {
	order := OrderNew{
		Market:    "ETH-EUR",
		Side:      SideBuy,
		OrderType: OrderTypeLimit,
	}

	bytes, err := json.Marshal(order)
	if err != nil {
		t.Error(err)
	}

	expected := "{\"market\":\"ETH-EUR\",\"side\":\"buy\",\"orderType\":\"limit\"}"
	actual := string(bytes)

	fmt.Printf("%q", actual)

	test.AssertEqual(t, expected, actual)
}

func TestMaxOrderUpdateMarshaller(t *testing.T) {
	order := OrderUpdate{
		Market:              "ETH-EUR",
		OrderId:             "123",
		Amount:              1.5,
		AmountQuote:         105.5,
		AmountRemaining:     10.5,
		Price:               2500.50,
		TriggerAmount:       10.2,
		TimeInForce:         TimeInForceDefault,
		SelfTradePrevention: SelfTradePreventionDefault,
		PostOnly:            true,
		ResponseRequired:    true,
	}

	bytes, err := json.Marshal(order)
	if err != nil {
		t.Error(err)
	}

	expected := "{\"market\":\"ETH-EUR\",\"orderId\":\"123\",\"amount\":1.5,\"amountQuote\":105.5,\"amountRemaining\":10.5,\"price\":2500.5,\"triggerAmount\":10.2,\"postOnly\":true,\"responseRequired\":true,\"timeInForce\":\"GTC\",\"selfTradePrevention\":\"decrementAndCancel\"}"
	actual := string(bytes)

	fmt.Printf("%q", actual)

	test.AssertEqual(t, expected, actual)
}

func TestMinOrderUpdateMarshaller(t *testing.T) {
	order := OrderUpdate{
		Market:  "ETH-EUR",
		OrderId: "123",
	}

	bytes, err := json.Marshal(order)
	if err != nil {
		t.Error(err)
	}

	expected := "{\"market\":\"ETH-EUR\",\"orderId\":\"123\"}"
	actual := string(bytes)

	fmt.Printf("%q", actual)

	test.AssertEqual(t, expected, actual)
}
