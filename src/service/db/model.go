package db

import (
	"time"
)

type CryptocurrencyInfo struct {
	ShortName string `json:"code"`
	FullName  string `json:"-"`
	Rate      string `json:"rate"`
	UnitPower int32  `json:"-"`
}
type CryptocurrencyInfoSlice []CryptocurrencyInfo

func (self CryptocurrencyInfoSlice) Len() int {
	return len(self)
}
func (self CryptocurrencyInfoSlice) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
func (self CryptocurrencyInfoSlice) Less(i, j int) bool {
	return self[i].ShortName < self[j].ShortName
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
	Seq             int64   `json:"seq"`
	UpdatedAt       int64   `json:"update_at"`
	TransactionId   string  `json:"txid"`
	DepositAmount   uint64  `json:"deposit_value"`
	BuyAmount       uint64  `json:"sent"`
	RateFloat       float64 `json:"rate"`
	Rate            string
	Height          uint64 `json:"height"`
	BuyAddr         string `json:"address"`
	CurrencyType    string `json:"coin_type"`
	DepositAddr     string `json:"deposit_address"`
}

type RewardRecord struct {
	Id            uint64    `json:"id"`
	Creation      time.Time `json:"rewardTime"`
	DepositId     uint64    `json:"-"`
	Address       string    `json:"address"`
	CalAmount     uint64    `json:"-"`
	CalAmountStr  string    `json:"rewardAmount"`
	SentAmount    uint64    `json:"-"`
	SentAmountStr string    `json:"sentAmount"`
	SentTime      Time      `json:"sentTime"`
	Sent          bool      `json:"sent"`
	RewardType    string    `json:"type"`
}

const (
	RewardBuyer            = "Buyer"
	RewardPromoter         = "Promoter"
	RewardSuperiorPromoter = "SuperiorPromoter"
)
