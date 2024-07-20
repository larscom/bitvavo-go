package bitvavo

import (
	"fmt"
	"net/url"
	"time"

	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/internal/util"
)

type DepositAsset struct {
	// The address to which cryptocurrencies can be sent to increase the account balance.
	//
	// NOTICE: for digital deposits
	Address string `json:"address"`

	// If a paymentid is supplied, attaching this to your deposit is required. This is mostly called a note, memo or tag.
	//
	// NOTICE: for digital deposits
	PaymentId string `json:"paymentid"`

	// IBAN number to wire your deposit to.
	//
	// NOTICE: for fiat deposits
	IBAN string `json:"iban"`

	// Optional code sometimes necessary for international transfers.
	//
	// NOTICE: for fiat deposits
	BIC string `json:"bic"`

	// Description which must be used for the deposit.
	//
	// NOTICE: for fiat deposits
	Description string `json:"description"`
}

type DepositHistoryParams struct {
	// When no symbol is specified, all deposits will be returned.
	Symbol string `json:"symbol"`

	// Return the limit most recent assets only.
	// Default: 500
	Limit uint64 `json:"limit"`

	// Return orders after start time.
	Start time.Time `json:"start"`

	// Return orders before end time.
	End time.Time `json:"end"`
}

func (d *DepositHistoryParams) Params() url.Values {
	params := make(url.Values)

	if d.Symbol != "" {
		params.Add("symbol", fmt.Sprint(d.Symbol))
	}
	if d.Limit > 0 {
		params.Add("limit", fmt.Sprint(d.Limit))
	}
	if !d.Start.IsZero() {
		params.Add("start", fmt.Sprint(d.Start.UnixMilli()))
	}
	if !d.End.IsZero() {
		params.Add("end", fmt.Sprint(d.End.UnixMilli()))
	}

	return params
}

type DepositHistory struct {
	// The time your deposit of symbol was received by Bitvavo.
	Timestamp int64 `json:"timestamp"`

	// The short name of the base currency you deposited with Bitvavo. For example, BTC for Bitcoin or EUR for euro.
	Symbol string `json:"symbol"`

	// The quantity of symbol you deposited with Bitvavo.
	Amount float64 `json:"amount"`

	// The identifier for the account you sent amount of symbol from. For example, NL89BANK0123456789 or a digital address (e.g: 14qViLJfdGaP4EeHnDyJbEGQysnCpwk3gd).
	Address string `json:"address"`

	// The identifier for this deposit. If you did not set an ID when you made this deposit, this parameter is not included in the response.
	//
	// NOTICE: digital currency only
	PaymentId string `json:"paymentId"`

	// The ID for this transaction on the blockchain.
	//
	// NOTICE: digital currency only
	TxId string `json:"txId"`

	// The transaction fee you paid to deposit amount of symbol on Bitvavo.
	Fee float64 `json:"fee"`

	// The current state of this deposit. Possible values are:
	// completed - amount of symbol has been added to your balance on Bitvavo.
	// canceled - this deposit could not be completed.
	//
	// NOTICE: fiat currency only
	Status string `json:"status"`
}

func (d *DepositHistory) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		timestamp = util.GetOrEmpty[float64]("timestamp", j)
		symbol    = util.GetOrEmpty[string]("symbol", j)
		amount    = util.GetOrEmpty[string]("amount", j)
		address   = util.GetOrEmpty[string]("address", j)
		paymentId = util.GetOrEmpty[string]("paymentId", j)
		txId      = util.GetOrEmpty[string]("txId", j)
		fee       = util.GetOrEmpty[string]("fee", j)
		status    = util.GetOrEmpty[string]("status", j)
	)

	d.Timestamp = int64(timestamp)
	d.Symbol = symbol
	d.Amount = util.IfOrElse(len(amount) > 0, func() float64 { return util.MustFloat64(amount) }, 0)
	d.Address = address
	d.PaymentId = paymentId
	d.TxId = txId
	d.Fee = util.IfOrElse(len(fee) > 0, func() float64 { return util.MustFloat64(fee) }, 0)
	d.Status = status

	return nil
}
