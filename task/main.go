package main

import (
	"fmt"
	"runtime/debug"

	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service"
	"github.com/spaco/affiliate/src/service/db"
	//	"github.com/spaco/affiliate/src/tracking_code"
	"encoding/json"
	"log"
	"os"

	"github.com/shopspring/decimal"
	client "github.com/spaco/affiliate/src/teller_client"
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
		logger.Println(string(debug.Stack()))
	}
}

func main() {
	defer deferFunc()
	db.OpenDb(&conf.Db)
	defer db.CloseDb()
	syncCryptocurrency()
	syncDeposit()
	sendReward()
	//	testSendReward()
	//	resp, err := client.Status("2g4AGfrjU91cQHAQhWZLnc9DmDDfmPVR8o9", "skycoin")
	//	if err == nil {
	//		fmt.Println(resp)
	//	}
	//	testProcessDeposit()
}

func syncCryptocurrency() {
	defer deferFunc()
	slice := client.Rate()
	service.SyncCryptocurrency(slice)
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

func testSendReward() {
	defer deferFunc()
	rrs := make([]db.RewardRecord, 0, 4)
	rrs = append(rrs, db.RewardRecord{Id: 1, Address: "2KSZSEoijudK6R6C4s7rJS9Qt1yfxvKfvao", SentAmount: 1000000})
	rrs = append(rrs, db.RewardRecord{Id: 2, Address: "2TNbiXocP6PxAD6rULiFkTHgkUoCXpjNttc", SentAmount: 2000000})
	client.SendCoin(rrs)
}

func testProcessDeposit() {
	jsonStr := `[{
		              "seq": 2,
		              "update_at": 1514726616,
		              "address": "2g4AGfrjU91cQHAQhWZLnc9DmDDfmPVR8o9",
		              "deposit_address": "2do3K1YLMy3Aq6EcPMdncEurP5BfAUdFPJj",
		              "txid": "fa92485d739e64e55f7a4beab9f5d7e6e23aa6f7e289bd5a5e7597f5e1fa4cf9",
		              "sent": 20000000,
		              "rate": 100,
		              "coin_type": "skycoin",
		              "deposit_value": 200000,
		              "height": 7456
		          }]`
	//	fmt.Println(jsonStr)
	drs := make([]db.DepositRecord, 0, 2)
	err := json.Unmarshal([]byte(jsonStr), &drs)
	if err != nil {
		panic(err)
	}
	for i, _ := range drs {
		drs[i].Rate = decimal.NewFromFloat(drs[i].RateFloat).String()
	}
	service.ProcessDeposit(drs, 3)
}
