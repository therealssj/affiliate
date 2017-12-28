package db

import (
	"time"
)

type CryptocurrencyInfo struct {
	ShortName string
	FullName  string
	Rate      string
}

type BuyAddrMapping struct {
	Id           uint64
	Creation     time.Time
	Address      string
	CurrencyType string
	DepositAddr  string
	Ref          string
}

type DepositRecord struct {
	Id              uint64
	Creation        time.Time
	MappingId       uint64
	RefAddr         string
	SuperiorRefAddr string
	Seq             int64  `json:"seq"`
	UpdatedAt       int64  `json:"update_at"`
	TransactionId   string `json:"txid"`
	DepositAmount   uint64 `json:"deposit_value"`
	BuyAmount       uint64 `json:"sent"`
	Rate            string `json:"rate"`
	Height          uint64 `json:"height"`
	BuyAddr         string `json:"address"`
	CurrencyType    string `json:"coin_type"`
	DepositAddr     string `json:"deposit_address"`
}

type RewardRecord struct {
	Id         uint64    `json:"id"`
	Creation   time.Time `json:"rewardTime"`
	DepositId  uint64    `json:"depositId"`
	Address    string    `json:"address"`
	CalAmount  uint64    `json:"rewardAmount"`
	SentAmount uint64    `json:"amount"`
	SentTime   time.Time `json:"sentTime"`
	Sent       bool      `json:"sent"`
	RewardType string    `json:"type"`
}

const (
	RewardBuyer            = "Buyer"
	RewardPromoter         = "Promoter"
	RewardSuperiorPromoter = "SuperiorPromoter"
)
