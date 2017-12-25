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
	Id             uint64
	Version        uint64
	Creation       time.Time
	LastModified   time.Time
	Address        string
	CurrencyType   string
	DepositAddr    string
	Ref            string
	DepositAmount  float64
	BuyAmount      uint64
	LastUpdated    time.Time
	TransactionIds string
	SentCoin       bool
}
