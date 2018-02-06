package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service"
	"github.com/spolabs/affiliate/src/service/db"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	//	"strings"
	"time"
)

var logger *log.Logger

func init() {
	os.MkdirAll(config.GetApiForTellerConfig().LogFolder, 0755)
	f, err := os.OpenFile(config.GetServerConfig().LogFolder+"api_for_teller.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "INFO", log.Ldate|log.Ltime)
}

func main() {
	//	drs := make([]db.DepositRecord, 0, 8)
	//	s := `[{"seq":11,"update_at":1517134128,"coin_type":"SKY","address":"25PN9qx8NKga2RFqHNv5xm9UBuowk5gi9pv","deposit_address":"LAWbVXeTL82vxjh21TNv6ALnMv2CT1mjL4","txid":"f37d9e96b84c5a7451993e5252da91d84c857a406caa4e22eb783e21eb8907a8","rate":"70","sent":119000000,"deposit_value":1700000,"height":12643}]`
	//	er := json.NewDecoder(strings.NewReader(s)).Decode(&drs)
	//	if er != nil {
	//		fmt.Println(er.Error())
	//	} else {
	//		fmt.Println(drs)
	//	}
	http.HandleFunc("/api/deposit/", depositHandler)
	http.HandleFunc("/api/reward/", rewardHandler)
	http.HandleFunc("/api/reward-status/", rewardStatusHandler)
	config := config.GetApiForTellerConfig()
	db.OpenDb(&config.Db)
	defer db.CloseDb()
	fmt.Printf("Listening on %s:%d", config.ListenIp, config.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.ListenIp, config.ListenPort), nil)
	if err != nil {
		println("ListenAndServe Error： %s", err)
	}
}

type JsonObj struct {
	Code   int8        `json:"code"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func depositHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverErr(w, r)
	if !chechAuthHeader(w, r) {
		return
	}
	drs := make([]db.DepositRecord, 0, 8)
	err := json.NewDecoder(r.Body).Decode(&drs)
	if !checkJsonErr(err, w, r) {
		return
	}
	service.ProcessDeposit(drs)
	success(w, r)
}

func rewardHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverErr(w, r)
	if !chechAuthHeader(w, r) {
		return
	}
	rrs := service.GetUnsentRewardRecord()
	if len(rrs) == 0 {
		json.NewEncoder(w).Encode(&JsonObj{Code: -1, ErrMsg: "no record"})
		return
	}
	slice := make([]rewardInfo, 0, len(rrs))
	for _, rr := range rrs {
		slice = append(slice, rewardInfo{rr.Id, rr.Address, rr.SentAmount})
	}
	json.NewEncoder(w).Encode(&JsonObj{0, "", slice})
}

type rewardInfo struct {
	Id         uint64 `json:"id"`
	Address    string `json:"address"`
	SentAmount uint64 `json:"amount"`
}

func rewardStatusHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverErr(w, r)
	if !chechAuthHeader(w, r) {
		return
	}
	ids := make([]uint64, 0, 8)
	err := json.NewDecoder(r.Body).Decode(&ids)
	if !checkJsonErr(err, w, r) {
		return
	}
	service.UpdateBatchRewardRecord(ids...)
	success(w, r)
}

func recoverErr(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		debug.PrintStack()
		logger.Println(debug.Stack())
		json.NewEncoder(w).Encode(&JsonObj{Code: 9, ErrMsg: fmt.Sprint(err)})
	}
}

func checkJsonErr(err error, w http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		json.NewEncoder(w).Encode(&JsonObj{Code: 1, ErrMsg: "cannot parse request to JSON:" + err.Error()})
		return false
	}
	return true
}
func chechAuthHeader(w http.ResponseWriter, r *http.Request) bool {
	//	if true {
	//		return true
	//	}
	conf := config.GetApiForTellerConfig()
	tsStr := r.Header.Get("timestamp")
	ts, err := strconv.Atoi(tsStr)
	if err != nil {
		json.NewEncoder(w).Encode(&JsonObj{Code: 2, ErrMsg: "invalid header timestamp: " + err.Error()})
		return false
	}
	timestamp := int64(ts)
	current := time.Now().Unix()
	if timestamp-current > 3 {
		json.NewEncoder(w).Encode(&JsonObj{Code: 3, ErrMsg: "client time error"})
		return false
	}
	if current-timestamp > int64(conf.AuthValidSec) {
		json.NewEncoder(w).Encode(&JsonObj{Code: 4, ErrMsg: "timestamp expired"})
		return false
	}
	hash := hmac.New(sha256.New, []byte(conf.AuthToken))
	hash.Write([]byte(tsStr))
	if hex.EncodeToString(hash.Sum(nil)) != r.Header.Get("auth") {
		json.NewEncoder(w).Encode(&JsonObj{Code: 5, ErrMsg: "auth verify error"})
		return false
	}
	return true
}

func success(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&JsonObj{})
}
