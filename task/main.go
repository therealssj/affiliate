package main

import (
	"fmt"
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service"
	"github.com/spaco/affiliate/src/service/db"
	//	"github.com/spaco/affiliate/src/tracking_code"
	client "github.com/spaco/affiliate/src/teller_client"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic Error: %s", err)
		}
	}()
	config := config.GetDaemonConfig()
	db.OpenDb(&config.Db)
	defer db.CloseDb()
	syncCryptocurrency()
	//	req := service.GetTellerReq()
}

func syncCryptocurrency() {
	currencyMap := service.AllCryptocurrencyMap()
	newCurrency := make([]*db.CryptocurrencyInfo, 0, 4)
	updateRateCur := make([]*db.CryptocurrencyInfo, 0, 4)
	for _, info := range client.Rate() {
		if old, ok := currencyMap[info.ShortName]; ok {
			if old.Rate != info.Rate {
				updateRateCur = append(updateRateCur, info)
			}
		} else {
			newCurrency = append(newCurrency, info)
		}
	}
	service.SyncCryptocurrency(newCurrency, updateRateCur)
}
