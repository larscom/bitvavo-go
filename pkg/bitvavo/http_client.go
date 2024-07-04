package bitvavo

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/larscom/bitvavo-go/internal/util"
)

const (
	bitvavoURL = "https://api.bitvavo.com/v2"

	headerRatelimit        = "Bitvavo-Ratelimit-Remaining"
	headerRatelimitResetAt = "Bitvavo-Ratelimit-Resetat"
	headerAccessKey        = "Bitvavo-Access-Key"
	headerAccessSignature  = "Bitvavo-Access-Signature"
	headerAccessTimestamp  = "Bitvavo-Access-Timestamp"
	headerAccessWindow     = "Bitvavo-Access-Window"
)

var emptyParams = make(url.Values)

type Params interface {
	Params() url.Values
}

type PublicAPI interface {
	// GetRateLimit returns the remaining rate limit.
	//
	// Default value: -1
	GetRateLimit() int64

	// GetRateLimitResetAt returns the time (local time) when the counter resets.
	GetRateLimitResetAt() time.Time

	// GetTime returns the current server time in milliseconds since 1 Jan 1970
	GetTime(ctx context.Context) (int64, error)

	// GetMarkets returns the available markets with their status (trading,halted,auction) and
	// available order types.
	GetMarkets(ctx context.Context) ([]Market, error)

	// GetMarkets returns the available markets with their status (trading,halted,auction) and
	// available order types for a single market (e.g: ETH-EUR)
	GetMarket(ctx context.Context, market string) (Market, error)

	// GetAssets returns information on the supported assets
	GetAssets(ctx context.Context) ([]Asset, error)

	// GetAsset returns information on the supported asset by symbol (e.g: ETH).
	GetAsset(ctx context.Context, symbol string) (Asset, error)

	// GetOrderBook returns a book with bids and asks for market.
	// That is, the buy and sell orders made by all Bitvavo users in a specific market (e.g: ETH-EUR).
	// The orders in the return parameters are sorted by price
	//
	// Optionally provide the depth (single value) to return the top depth orders only.
	GetOrderBook(ctx context.Context, market string, depth ...uint64) (Book, error)

	// GetTrades returns the list of all trades made by all Bitvavo users for market (e.g: ETH-EUR).
	// That is, the trades that have been executed in the past.
	//
	// Optionally provide extra params (see: TradeParams)
	GetTrades(ctx context.Context, market string, params ...Params) ([]Trade, error)

	// GetCandles returns the Open, High, Low, Close, Volume (OHLCV) data you use to create candlestick charts
	// for market with interval time between each candlestick (e.g: market=ETH-EUR interval=5m)
	//
	// Optionally provide extra params (see: CandleParams)
	GetCandles(ctx context.Context, market string, interval Interval, params ...Params) ([]CandleOnly, error)

	// GetTickerPrices returns price of the latest trades on Bitvavo for all markets.
	GetTickerPrices(ctx context.Context) ([]TickerPrice, error)

	// GetTickerPrice returns price of the latest trades on Bitvavo for a single market (e.g: ETH-EUR).
	GetTickerPrice(ctx context.Context, market string) (TickerPrice, error)

	// GetTickerBooks returns the highest buy and the lowest sell prices currently available for
	// all markets in the Bitvavo order book.
	GetTickerBooks(ctx context.Context) ([]TickerBook, error)

	// GetTickerBook returns the highest buy and the lowest sell prices currently
	// available for a single market (e.g: ETH-EUR) in the Bitvavo order book.
	GetTickerBook(ctx context.Context, market string) (TickerBook, error)

	// GetTickers24h returns high, low, open, last, and volume information for trades and orders for all markets over the previous 24 hours.
	GetTickers24h(ctx context.Context) ([]Ticker24hData, error)

	// GetTicker24h returns high, low, open, last, and volume information for trades and orders for a single market over the previous 24 hours.
	GetTicker24h(ctx context.Context, market string) (Ticker24hData, error)
}

type PrivateAPI interface {
	PublicAPI

	// GetBalance returns the balance on the account.
	// Optionally provide the symbol to filter for in uppercase (e.g: ETH)
	GetBalance(ctx context.Context, symbol ...string) ([]Balance, error)

	// GetAccount returns trading volume and fees for account.
	GetAccount(ctx context.Context) (Account, error)

	// GetTrades returns historic trades for your account for market (e.g: ETH-EUR)
	//
	// Optionally provide extra params (see: TradeParams)
	GetTradesHistoric(ctx context.Context, market string, params ...Params) ([]TradeHistoric, error)

	// GetOrders returns data for multiple orders at once for market (e.g: ETH-EUR)
	//
	// Optionally provide extra params (see: OrderParams)
	GetOrders(ctx context.Context, market string, params ...Params) ([]Order, error)

	// GetOrdersOpen returns all open orders for market (e.g: ETH-EUR) or all open orders
	// if no market is given.
	GetOrdersOpen(ctx context.Context, market ...string) ([]Order, error)

	// GetOrder returns the order by market and ID
	GetOrder(ctx context.Context, market string, orderId string) (Order, error)

	// CancelOrders cancels multiple orders at once.
	// Either for an entire market (e.g: ETH-EUR) or for the entire account if you
	// omit the market.
	//
	// It returns a slice of orderId's of which are canceled
	CancelOrders(ctx context.Context, market ...string) ([]string, error)

	// CancelOrder cancels a single order by ID for the specific market (e.g: ETH-EUR)
	//
	// It returns the canceled orderId if it was canceled
	CancelOrder(ctx context.Context, market string, orderId string) (string, error)

	// NewOrder places a new order on the exchange.
	//
	// It returns the new order if it was successfully created
	NewOrder(ctx context.Context, market string, side string, orderType string, order OrderNew) (Order, error)

	// UpdateOrder updates an existing order on the exchange.
	//
	// It returns the updated order if it was successfully updated
	UpdateOrder(ctx context.Context, market string, orderId string, order OrderUpdate) (Order, error)

	// GetDepositAsset returns deposit address (with paymentid for some assets)
	// or bank account information to increase your balance for a specific symbol (e.g: ETH)
	GetDepositAsset(ctx context.Context, symbol string) (DepositAsset, error)

	// GetDepositHistory returns the deposit history of the account.
	//
	// Optionally provide extra params (see: DepositHistoryParams)
	GetDepositHistory(ctx context.Context, params ...Params) ([]DepositHistory, error)

	// GetWithdrawalHistory returns the withdrawal history of the account.
	//
	// Optionally provide extra params (see: WithdrawalHistoryParams)
	GetWithdrawalHistory(ctx context.Context, params ...Params) ([]WithdrawalHistory, error)

	// Withdraw requests a withdrawal to an external cryptocurrency address or verified bank account.
	// Please note that 2FA and address confirmation by e-mail are disabled for API withdrawals.
	Withdraw(ctx context.Context, symbol string, amount float64, address string, withdrawal Withdrawal) (WithDrawalResponse, error)
}

type privateConfig struct {
	apiKey     string
	apiSecret  string
	windowTime uint16
}

type httpClient struct {
	mu               sync.RWMutex
	ratelimit        int64
	ratelimitResetAt time.Time

	config *privateConfig
}

func NewPublicHTTPClient() PublicAPI {
	return &httpClient{
		ratelimit: -1,
	}
}

func (c *httpClient) GetRateLimit() int64 {
	return c.ratelimit
}

func (c *httpClient) GetRateLimitResetAt() time.Time {
	return c.ratelimitResetAt
}

func (c *httpClient) GetTime(ctx context.Context) (int64, error) {
	resp, err := httpGet[map[string]float64](
		ctx,
		fmt.Sprintf("%s/time", bitvavoURL),
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
	if err != nil {
		return 0, err
	}

	return int64(resp["time"]), nil
}

func (c *httpClient) GetMarkets(ctx context.Context) ([]Market, error) {
	return httpGet[[]Market](
		ctx,
		fmt.Sprintf("%s/markets", bitvavoURL),
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetMarket(ctx context.Context, market string) (Market, error) {
	params := make(url.Values)
	params.Add("market", market)

	return httpGet[Market](
		ctx,
		fmt.Sprintf("%s/markets", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetAssets(ctx context.Context) ([]Asset, error) {
	return httpGet[[]Asset](
		ctx,
		fmt.Sprintf("%s/assets", bitvavoURL),
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetAsset(ctx context.Context, symbol string) (Asset, error) {
	params := make(url.Values)
	params.Add("symbol", symbol)

	return httpGet[Asset](
		ctx,
		fmt.Sprintf("%s/assets", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetOrderBook(ctx context.Context, market string, depth ...uint64) (Book, error) {
	params := make(url.Values)
	if len(depth) > 0 {
		params.Add("depth", fmt.Sprint(depth[0]))
	}

	return httpGet[Book](
		ctx,
		fmt.Sprintf("%s/%s/book", bitvavoURL, market),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetTrades(ctx context.Context, market string, opt ...Params) ([]Trade, error) {
	params := make(url.Values)
	if len(opt) > 0 {
		params = opt[0].Params()
	}
	return httpGet[[]Trade](
		ctx,
		fmt.Sprintf("%s/%s/trades", bitvavoURL, market),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetCandles(ctx context.Context, market string, interval Interval, opt ...Params) ([]CandleOnly, error) {
	params := make(url.Values)
	if len(opt) > 0 {
		params = opt[0].Params()
	}
	params.Add("interval", interval.Value)

	return httpGet[[]CandleOnly](
		ctx,
		fmt.Sprintf("%s/%s/candles", bitvavoURL, market),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetTickerPrices(ctx context.Context) ([]TickerPrice, error) {
	return httpGet[[]TickerPrice](
		ctx,
		fmt.Sprintf("%s/ticker/price", bitvavoURL),
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetTickerPrice(ctx context.Context, market string) (TickerPrice, error) {
	params := make(url.Values)
	params.Add("market", market)

	return httpGet[TickerPrice](
		ctx,
		fmt.Sprintf("%s/ticker/price", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetTickerBooks(ctx context.Context) ([]TickerBook, error) {
	return httpGet[[]TickerBook](
		ctx,
		fmt.Sprintf("%s/ticker/book", bitvavoURL),
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetTickerBook(ctx context.Context, market string) (TickerBook, error) {
	params := make(url.Values)
	params.Add("market", market)

	return httpGet[TickerBook](
		ctx,
		fmt.Sprintf("%s/ticker/book", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetTickers24h(ctx context.Context) ([]Ticker24hData, error) {
	return httpGet[[]Ticker24hData](
		ctx,
		fmt.Sprintf("%s/ticker/24h", bitvavoURL),
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) GetTicker24h(ctx context.Context, market string) (Ticker24hData, error) {
	params := make(url.Values)
	params.Add("market", market)

	return httpGet[Ticker24hData](
		ctx,
		fmt.Sprintf("%s/ticker/24h", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		nil,
	)
}

func (c *httpClient) updateRateLimit(ratelimit int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ratelimit = ratelimit
}

func (c *httpClient) updateRateLimitResetAt(resetAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ratelimitResetAt = resetAt
}

func NewPrivateHTTPClient(apiKey, apiSecret string, windowTimeMs ...uint16) PrivateAPI {
	windowTime := util.IfOrElse(len(windowTimeMs) > 0, func() uint16 { return windowTimeMs[0] }, 10000)
	if windowTime == 0 || windowTime > 60000 {
		panic("windowTimeMs must be > 0 and <= 60000")
	}

	return &httpClient{
		ratelimit: -1,
		config: &privateConfig{
			apiKey:     apiKey,
			apiSecret:  apiSecret,
			windowTime: windowTime,
		},
	}
}

func (c *httpClient) GetBalance(ctx context.Context, symbol ...string) ([]Balance, error) {
	params := make(url.Values)
	if len(symbol) > 0 {
		params.Add("symbol", symbol[0])
	}

	return httpGet[[]Balance](
		ctx,
		fmt.Sprintf("%s/balance", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) GetAccount(ctx context.Context) (Account, error) {
	return httpGet[Account](
		ctx,
		fmt.Sprintf("%s/account", bitvavoURL),
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) GetOrders(ctx context.Context, market string, opt ...Params) ([]Order, error) {
	params := make(url.Values)
	if len(opt) > 0 {
		params = opt[0].Params()
	}
	params.Add("market", market)

	return httpGet[[]Order](
		ctx,
		fmt.Sprintf("%s/orders", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) GetOrdersOpen(ctx context.Context, market ...string) ([]Order, error) {
	params := make(url.Values)
	if len(market) > 0 {
		params.Add("market", market[0])
	}

	return httpGet[[]Order](
		ctx,
		fmt.Sprintf("%s/ordersOpen", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) GetOrder(ctx context.Context, market string, orderId string) (Order, error) {
	params := make(url.Values)
	params.Add("market", market)
	params.Add("orderId", orderId)

	return httpGet[Order](
		ctx,
		fmt.Sprintf("%s/order", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) CancelOrders(ctx context.Context, market ...string) ([]string, error) {
	params := make(url.Values)
	if len(market) > 0 {
		params.Add("market", market[0])
	}

	resp, err := httpDelete[[]map[string]string](
		ctx,
		fmt.Sprintf("%s/orders", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
	if err != nil {
		return nil, err
	}

	orderIds := make([]string, len(resp))
	for i := 0; i < len(orderIds); i++ {
		orderIds[i] = resp[i]["orderId"]
	}

	return orderIds, nil
}

func (c *httpClient) CancelOrder(ctx context.Context, market string, orderId string) (string, error) {
	params := make(url.Values)
	params.Add("market", market)
	params.Add("orderId", orderId)

	resp, err := httpDelete[map[string]string](
		ctx,
		fmt.Sprintf("%s/order", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
	if err != nil {
		return "", err
	}

	return resp["orderId"], nil
}

func (c *httpClient) NewOrder(ctx context.Context, market string, side string, orderType string, order OrderNew) (Order, error) {
	order.Market = market
	order.Side = side
	order.OrderType = orderType
	return httpPost[Order](
		ctx,
		fmt.Sprintf("%s/order", bitvavoURL),
		order,
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) UpdateOrder(ctx context.Context, market string, orderId string, order OrderUpdate) (Order, error) {
	order.Market = market
	order.OrderId = orderId

	return httpPut[Order](
		ctx,
		fmt.Sprintf("%s/order", bitvavoURL),
		order,
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) GetTradesHistoric(ctx context.Context, market string, opt ...Params) ([]TradeHistoric, error) {
	params := make(url.Values)
	if len(opt) > 0 {
		params = opt[0].Params()
	}
	params.Add("market", market)

	return httpGet[[]TradeHistoric](
		ctx,
		fmt.Sprintf("%s/trades", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) GetDepositAsset(ctx context.Context, symbol string) (DepositAsset, error) {
	params := make(url.Values)
	params.Add("symbol", symbol)

	return httpGet[DepositAsset](
		ctx,
		fmt.Sprintf("%s/deposit", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) GetDepositHistory(ctx context.Context, opt ...Params) ([]DepositHistory, error) {
	params := make(url.Values)
	if len(opt) > 0 {
		params = opt[0].Params()
	}
	return httpGet[[]DepositHistory](
		ctx,
		fmt.Sprintf("%s/depositHistory", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) GetWithdrawalHistory(ctx context.Context, opt ...Params) ([]WithdrawalHistory, error) {
	params := make(url.Values)
	if len(opt) > 0 {
		params = opt[0].Params()
	}
	return httpGet[[]WithdrawalHistory](
		ctx,
		fmt.Sprintf("%s/withdrawalHistory", bitvavoURL),
		params,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}

func (c *httpClient) Withdraw(ctx context.Context, symbol string, amount float64, address string, withdrawal Withdrawal) (WithDrawalResponse, error) {
	withdrawal.Symbol = symbol
	withdrawal.Amount = amount
	withdrawal.Address = address

	return httpPost[WithDrawalResponse](
		ctx,
		fmt.Sprintf("%s/withdrawal", bitvavoURL),
		withdrawal,
		emptyParams,
		c.updateRateLimit,
		c.updateRateLimitResetAt,
		c.config,
	)
}
