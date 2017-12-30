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

var conf *config.DaemonConfig

func init() {
	conf = config.GetDaemonConfig()
	os.MkdirAll(conf.LogFolder, 0755)
	f, err := os.OpenFile(conf.LogFolder+"task.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	logger = log.New(f, "INFO", log.Ldate|log.Ltime)
}

func deferFunc() {
	if err := recover(); err != nil {
		fmt.Printf("Panic Error: %s", err)
		debug.PrintStack()
		logger.Println(debug.Stack())
	}
}

func main() {
	defer deferFunc()
	db.OpenDb(&conf.Db)
	defer db.CloseDb()
	syncCryptocurrency()
	syncDeposit()
	sendReward()
}

func syncCryptocurrency() {
	defer deferFunc()
	currencyMap := service.AllCryptocurrencyMap()
	newCurrency := make([]db.CryptocurrencyInfo, 0, 4)
	updateRateCur := make([]db.CryptocurrencyInfo, 0, 4)
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
	defer deferFunc()
	req := service.GetTellerReq()
	for {
		depositResp := client.Deposite(req)
		req = depositResp.NextSeq
		if len(depositResp.Deposits) == 0 {
			service.SaveTellerReq(req)
		} else {
			service.ProcessDeposit(depositResp.Deposits, req)
		}
		if !depositResp.GoOn {
			break
		}
	}
}

func sendReward() {
	defer deferFunc()
	rrs := service.GetUnsentRewardRecord()
	if len(rrs) == 0 {
		return
	}
	ids := make([]uint64, len(rrs))
	for _, rr := range rrs {
		ids = append(ids, rr.Id)
	}
	client.SendCoin(rrs)
	service.UpdateBatchRewardRecord(ids...)
}
