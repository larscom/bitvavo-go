package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/larscom/bitvavo-go/pkg/bitvavo"
)

func main() {
	_ = godotenv.Load()

	var (
		apiKey    = os.Getenv("API_KEY")
		apiSecret = os.Getenv("API_SECRET")
	)

	client := bitvavo.NewPrivateHTTPClient(apiKey, apiSecret)
	markets, err := client.GetMarkets(context.Background())
	if err != nil {
		panic(err)
	}

	tradingMarkets := make([]string, 0)
	for _, market := range markets {
		if market.Status == bitvavo.MarketStatusTrading {
			tradingMarkets = append(tradingMarkets, market.Market)
		}
	}

	log.Printf("Subscribing to %d markets\n", len(tradingMarkets))
	<-time.After(time.Second * 2)

	listener := bitvavo.NewTickerListener()
	defer func() {
		_ = listener.Close()
	}()

	chn, err := listener.Subscribe(tradingMarkets)
	if err != nil {
		panic(err)
	}

	for e := range chn {
		if e.Error != nil {
			panic(e.Error)
		}
		log.Println(e.Value)
	}
}
