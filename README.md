# BITVAVO-GO

[![Go Report Card](https://goreportcard.com/badge/github.com/larscom/bitvavo-go)](https://goreportcard.com/report/github.com/larscom/bitvavo-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/larscom/bitvavo-go.svg)](https://pkg.go.dev/github.com/larscom/bitvavo-go)

> GO **thread safe** library (WebSockets / HTTP) for Bitvavo v2 API (see: https://docs.bitvavo.com)

Listen to all events occurring on the Bitvavo platform (tickers, tickers24h, candles, books, trades, orders, fills) using websockets. With the HTTP client you can do things like placing orders or withdraw assets from your account.

## 📒 Features

- WebSocket Listeners -- Read only
  - Book
  - Candles
  - Trades
  - Ticker
  - Ticker 24h
  - Orders/Fills
- HTTP Client -- Read / Write
  - Market data endpoints
  - Account endpoints
  - Synchronization endpoints
  - Trading endpoints
  - Transfer endpoints

## 🚀 Installation

```shell
go get github.com/larscom/bitvavo-go@latest
```

## 💡 Usage

```shell
import "github.com/larscom/bitvavo-go/pkg/bitvavo"
```

## 👂 WebSocket

For each event on the Bitvavo platform there is a listener available. A listener wraps a websocket connection, you can also implement your own wrapper arround the websocket. The listeners handle everything for you, like resubscribing and reauthenticating when the connection has been lost.

### Public listeners

- BookListener
- CandlesListener
- TickerListener
- Ticker24hListener
- TradesListener

```go
import "github.com/larscom/bitvavo-go/pkg/bitvavo"

func main() {
	// listen for candle (public) events
	listener := bitvavo.NewCandlesListener()
	defer listener.Close()

	chn, err := listener.Listen([]string{"ETH-EUR"}, []bitvavo.Interval{bitvavo.INTERVAL_1M})
	if err != nil {
		panic(err)
	}

	for event := range chn {
		if event.Error != nil {
			panic(event.Error)
		}
		log.Println(event.Value)
	}
}

```

### Private listeners

- OrderListener
- FillListener

```go
import "github.com/larscom/bitvavo-go/pkg/bitvavo"

func main() {
	// listen for order (private) events
	listener := bitvavo.NewOrderListener("MY_API_KEY", "MY_API_SECRET")
	defer listener.Close()

	chn, err := listener.Listen([]string{"ETH-EUR"})
	if err != nil {
		panic(err)
	}

	for event := range chn {
		if event.Error != nil {
			panic(event.Error)
		}
		log.Println(event.Value)
	}
}

```

### Create custom listener

It's possible to create your own wrapper arround the websocket and listen to multiple events at the same time.

```go
import "github.com/larscom/bitvavo-go/pkg/bitvavo"

func main() {
	onMessage := func(data bitvavo.WebSocketEventData, err error) {
			if err != nil {
				// oh no error!
			} else if data.Event == bitvavo.EVENT_BOOK {
				// decode into Book
				var book bitvavo.Book
				data.Decode(&book)
			} else if data.Event == bitvavo.EVENT_CANDLE {
				// decode into Candle
				var candle bitvavo.Candle
				data.Decode(&candle)
			}
			// etc
		}

		onReconnect := func() {
			// gets called when successfully reconnected
		}

		ws, err := bitvavo.NewWebSocket(context.Background(), onMessage, onReconnect)
		// do stuff with ws
}
```

## 🌐 HTTP

The HTTP client implements 2 interfaces (PrivateAPI and PublicAPI)

If you need both private and public endpoints you can create a private http client as it includes both public and private endpoints.

### Private and Public endpoints

```go
import "github.com/larscom/bitvavo-go/pkg/bitvavo"

func main() {
	// private http client (includes public as well)
	client := bitvavo.NewPrivateHTTPClient("MY_API_KEY", "MY_API_SECRET")

	orders, err := client.GetOrders(context.Background(), "ETH-EUR")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Orders", orders)
}

```

### Public endpoints only

```go
import "github.com/larscom/bitvavo-go/pkg/bitvavo"

func main() {
	// public http client
	client := bitvavo.NewPublicHTTPClient()

	markets, err := client.GetMarkets(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Markets", markets)
}

```

### Endpoints with params

Some endpoints have additional params which you can provide.

```go
import "github.com/larscom/bitvavo-go/pkg/bitvavo"

func main() {
	client := bitvavo.NewPublicHTTPClient()

	// limit to 100 trades
	params := &bitvavo.TradeParams{
		Limit: 100,
	}
	trades, err := client.GetTrades(context.Background(), "ETH-EUR", params)
}

```

## 👉🏼 Run example

There is an example that uses the ticker listener for ticker events
and HTTP client to retrieve the trading markets.

You can run this example by cloning this project and running:

`make run` or without make: `go run ./example/main.go`

This command will subscribe to all available trading markets and log the received tickers.

### Private

If you want to test private endpoints and listeners you can place a `.env` file in the root of the project.

> .env file

```shell
API_KEY=MY_API_KEY
API_SECRET=MY_API_SECRET
```
