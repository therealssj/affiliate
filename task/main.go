package main

import (
	"fmt"
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service"
	"github.com/spaco/affiliate/src/service/db"
	"runtime/debug"
	//	"github.com/spaco/affiliate/src/tracking_code"
	client "github.com/spaco/affiliate/src/teller_client"
	"log"
	"os"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var logger *log.Logger

func init() {
	os.MkdirAll(config.GetDaemonConfig().LogFolder, 0755)
	f, err := os.OpenFile(config.GetDaemonConfig().LogFolder+"task.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	logger = log.New(f, "INFO", log.Ldate|log.Ltime)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic Error: %s", err)
			debug.PrintStack()
			logger.Println(debug.Stack())
		}
	}()
	config := config.GetDaemonConfig()
	db.OpenDb(&config.Db)
	defer db.CloseDb()
	syncCryptocurrency()
	syncDeposit()
}

func syncCryptocurrency() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic Error: %s", err)
			debug.PrintStack()
			logger.Println(debug.Stack())
		}
	}()
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

func syncDeposit() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic Error: %s", err)
			debug.PrintStack()
			logger.Println(debug.Stack())
		}
	}()
	req := service.GetTellerReq()
	for {
		depositResp := client.Deposite(req)
		req = depositResp.NextSeq
		service.ProcessDeposit(depositResp.Deposit, req)
		if !depositResp.GoOn {
			break
		}
	}
}
