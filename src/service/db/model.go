package db

import (
	"time"
)

type CryptocurrencyInfo struct {
	ShortName string
	FullName  string
	Rate      float32
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
	Id            uint64
	Creation      time.Time
	MappingId     uint64
	Seq           int64
	UpdatedAt     uint64
	TransactionId string `json:"Txid"`
	DepositAmount float32
	BuyAmount     uint64
	Rate          float32
	Height        uint64
}

type AwardRecord struct {
	Id        uint64
	Creation  time.Time
	DepositId uint64
}
