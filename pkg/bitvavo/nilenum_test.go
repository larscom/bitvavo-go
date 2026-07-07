package bitvavo

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/v2/internal/test"
)

// Every case here panicked before enum.Parse's nil result was guarded.

func TestOrderUnmarshalMissingSelfTradePrevention(t *testing.T) {
	body := []byte(`{"orderId":"1","market":"BTC-EUR","status":"filled","side":"buy","orderType":"market","fills":[]}`)

	var o Order
	if err := json.Unmarshal(body, &o); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, SelfTradePrevention{}, o.SelfTradePrevention)
}

func TestOrderUnmarshalWithFills(t *testing.T) {
	body := []byte(`{"orderId":"1","market":"BTC-EUR","status":"filled","side":"buy","orderType":"market","selfTradePrevention":"decrementAndCancel","fills":[{"id":"f1","timestamp":1,"amount":"1","price":"1","taker":true,"fee":"0","feeCurrency":"EUR","settled":true}]}`)

	var o Order
	if err := json.Unmarshal(body, &o); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, 1, len(o.Fills))
	test.AssertEqual(t, "f1", o.Fills[0].Id)
}

func TestCandleUnmarshal12hInterval(t *testing.T) {
	body := []byte(`{"market":"BTC-EUR","interval":"12h","candle":[[1,"1","1","1","1","1"]]}`)

	var c Candle
	if err := json.Unmarshal(body, &c); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, Interval12h, c.Interval)
}

func TestCandleUnmarshal1dInterval(t *testing.T) {
	body := []byte(`{"market":"BTC-EUR","interval":"1d","candle":[[1,"1","1","1","1","1"]]}`)

	var c Candle
	if err := json.Unmarshal(body, &c); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, Interval1d, c.Interval)
}

func TestFillUnmarshalMissingSide(t *testing.T) {
	body := []byte(`{"fillId":"f1","market":"BTC-EUR","orderId":"o1","timestamp":1,"amount":"1","price":"1","taker":true,"fee":"0","feeCurrency":"EUR","settled":true}`)

	var f Fill
	if err := json.Unmarshal(body, &f); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, Side{}, f.Side)
}

func TestTradeUnmarshalUnrecognizedSide(t *testing.T) {
	body := []byte(`{"id":"t1","market":"BTC-EUR","amount":"1","price":"1","side":"unknown","timestamp":1}`)

	var tr Trade
	if err := json.Unmarshal(body, &tr); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, Side{}, tr.Side)
}

func TestAssetUnmarshalUnrecognizedStatus(t *testing.T) {
	body := []byte(`{"symbol":"BTC","name":"Bitcoin","decimals":8,"depositStatus":"UNKNOWN","withdrawalStatus":"UNKNOWN"}`)

	var a Asset
	if err := json.Unmarshal(body, &a); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, DepositStatus{}, a.DepositStatus)
	test.AssertEqual(t, WithdrawalStatus{}, a.WithdrawalStatus)
}

func TestMarketUnmarshalUnrecognizedStatus(t *testing.T) {
	body := []byte(`{"market":"BTC-EUR","status":"unknown","base":"BTC","quote":"EUR"}`)

	var m Market
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, MarketStatus{}, m.Status)
}

func TestWithdrawalHistoryUnmarshalUnrecognizedStatus(t *testing.T) {
	body := []byte(`{"timestamp":1,"symbol":"BTC","amount":"1","address":"a","fee":"0","status":"unknown"}`)

	var w WithdrawalHistory
	if err := json.Unmarshal(body, &w); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, WithdrawalHistoryStatus{}, w.Status)
}

func TestSubscribedUnmarshalUnrecognizedChannel(t *testing.T) {
	body := []byte(`{"subscriptions":{"unknownChannel":["BTC-EUR"]}}`)

	var s Subscribed
	if err := json.Unmarshal(body, &s); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, 1, len(s.Subscriptions[Channel{}]))
}

func TestWebSocketEventDataUnmarshalUnrecognizedEvent(t *testing.T) {
	body := []byte(`{"event":"somethingNew"}`)

	var d WebSocketEventData
	if err := json.Unmarshal(body, &d); err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, WebSocketEvent{}, d.Event)
}
