package bitvavo

import (
	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/v2/internal/util"
)

type Book struct {
	// The market which was requested in the subscription.
	Market string `json:"market"`

	// Integer which is increased by one for every update to the book. Useful for synchronizing. Resets to zero after restarting the matching engine.
	Nonce int64 `json:"nonce"`

	// Slice with all bids in the format [price, size], where a size of 0 means orders are no longer present at that price level,
	// otherwise the returned size is the new total size on that price level.
	Bids []Page `json:"bids"`

	// Slice with all asks in the format [price, size], where a size of 0 means orders are no longer present at that price level,
	// otherwise the returned size is the new total size on that price level.
	Asks []Page `json:"asks"`
}

type Page struct {
	// Bid / ask price.
	Price string `json:"price"`

	//  Size of "0" means orders are no longer present at that price level, otherwise the returned size is the new total size on that price level.
	Size string `json:"size"`
}

func (b *Book) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		market    = util.GetOrEmpty[string]("market", j)
		nonce     = util.GetOrEmpty[float64]("nonce", j)
		bidEvents = util.GetOrEmpty[[]any]("bids", j)
		askEvents = util.GetOrEmpty[[]any]("asks", j)
	)

	bids := make([]Page, len(bidEvents))
	for i := 0; i < len(bidEvents); i++ {
		price := bidEvents[i].([]any)[0].(string)
		size := bidEvents[i].([]any)[1].(string)

		bids[i] = Page{
			Price: price,
			Size:  size,
		}
	}

	asks := make([]Page, len(askEvents))
	for i := 0; i < len(askEvents); i++ {
		price := askEvents[i].([]any)[0].(string)
		size := askEvents[i].([]any)[1].(string)

		asks[i] = Page{
			Price: price,
			Size:  size,
		}
	}

	b.Market = market
	b.Nonce = int64(nonce)
	b.Bids = bids
	b.Asks = asks

	return nil
}
