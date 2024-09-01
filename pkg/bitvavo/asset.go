package bitvavo

import (
	"github.com/goccy/go-json"
	"github.com/larscom/bitvavo-go/v2/internal/util"
	"github.com/orsinium-labs/enum"
)

type DepositStatus enum.Member[string]

var (
	depositStatus        = enum.NewBuilder[string, DepositStatus]()
	DepositStatusTrading = depositStatus.Add(DepositStatus{"OK"})
	DepositStatusHalted  = depositStatus.Add(DepositStatus{"MAINTENANCE"})
	DepositStatusAuction = depositStatus.Add(DepositStatus{"DELISTED"})
	depositStatuses      = depositStatus.Enum()
)

type WithdrawalStatus enum.Member[string]

var (
	withdrawalStatus        = enum.NewBuilder[string, WithdrawalStatus]()
	WithdrawalStatusTrading = withdrawalStatus.Add(WithdrawalStatus{"OK"})
	WithdrawalStatusHalted  = withdrawalStatus.Add(WithdrawalStatus{"MAINTENANCE"})
	WithdrawalStatusAuction = withdrawalStatus.Add(WithdrawalStatus{"DELISTED"})
	withdrawalStatuses      = withdrawalStatus.Enum()
)

type Asset struct {
	// Short version of the asset name used in market names.
	Symbol string `json:"symbol"`

	// The full name of the asset.
	Name string `json:"name"`

	// The precision used for specifying amounts.
	Decimals int64 `json:"decimals"`

	// Fixed fee for depositing this asset.
	DepositFee string `json:"depositFee"`

	// The minimum amount of network confirmations required before this asset is credited to your account.
	DepositConfirmations int64 `json:"depositConfirmations"`

	// The current deposit status.
	DepositStatus DepositStatus `json:"depositStatus"`

	// Fixed fee for withdrawing this asset.
	WithdrawalFee string `json:"withdrawalFee"`

	// The minimum amount for which a withdrawal can be made.
	WithdrawalMinAmount string `json:"withdrawalMinAmount"`

	// The current withdrawal status.
	WithdrawalStatus WithdrawalStatus `json:"withdrawalStatus"`

	// Supported networks.
	Networks []string `json:"networks"`

	// Shows the reason if withdrawalStatus or depositStatus is not OK.
	Message string `json:"message"`
}

func (m *Asset) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		symbol               = util.GetOrEmpty[string]("symbol", j)
		name                 = util.GetOrEmpty[string]("name", j)
		decimals             = util.GetOrEmpty[float64]("decimals", j)
		depositFee           = util.GetOrEmpty[string]("depositFee", j)
		depositConfirmations = util.GetOrEmpty[float64]("depositConfirmations", j)
		depositStatus        = util.GetOrEmpty[string]("depositStatus", j)
		withdrawalFee        = util.GetOrEmpty[string]("withdrawalFee", j)
		withdrawalMinAmount  = util.GetOrEmpty[string]("withdrawalMinAmount", j)
		withdrawalStatus     = util.GetOrEmpty[string]("withdrawalStatus", j)
		networksAny          = util.GetOrEmpty[[]any]("networks", j)
		message              = util.GetOrEmpty[string]("message", j)
	)

	networks := make([]string, len(networksAny))
	for i := 0; i < len(networksAny); i++ {
		networks[i] = networksAny[i].(string)
	}

	m.Symbol = symbol
	m.Name = name
	m.Decimals = int64(decimals)
	m.DepositFee = depositFee
	m.DepositConfirmations = int64(depositConfirmations)
	m.DepositStatus = *depositStatuses.Parse(depositStatus)
	m.WithdrawalFee = withdrawalFee
	m.WithdrawalMinAmount = withdrawalMinAmount
	m.WithdrawalStatus = *withdrawalStatuses.Parse(withdrawalStatus)
	m.Networks = networks
	m.Message = message

	return nil
}
