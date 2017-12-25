package teller_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
)

type jsonResp struct {
	Code   uint8           `json:"code"`
	ErrMsg string          `json:"errmsg"`
	Data   json.RawMessage `json:"data"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func httpGet(url string) []byte {
	config := config.GetServerConfig()
	resp, err := http.Get(config.Teller.ContextPath + url)
	checkErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	return body
}
func httpPost(url string, json []byte) []byte {
	config := config.GetServerConfig()
	resp, err := http.Post(config.Teller.ContextPath+url,
		"application/json",
		bytes.NewBuffer(json))
	checkErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	return body
}

type bindResp struct {
	Address      string `json:"address"`
	CurrencyType string `json:"coin_type"`
}

func Bind(currencyType string, address string) string {
	resp := httpPost("/api/bind/", []byte(fmt.Sprintf(`{"address":"%s","coin_type":"%s"}`, address, currencyType)))
	jsonObj := new(jsonResp)
	err := json.Unmarshal(resp, &jsonObj)
	checkErr(err)
	if jsonObj.Code != 0 {
		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	bResp := new(bindResp)
	err = json.Unmarshal(jsonObj.Data, &bResp)
	checkErr(err)
	if bResp.CurrencyType != currencyType {
		panic("return not same currency type")
	}
	return randStringRunes(32)
}

type rateResp struct {
	TokenType string     `json:"tokenType"`
	Rate      float32    `json:"rate"`
	AllCoin   []coinResp `json:"allcoin"`
}

type coinResp struct {
	Name string `json:"coin_name"`
	Code string `json:"coin_code"`
	Rate string `json:"coin_rate"`
}

func Rate() []*db.CryptocurrencyInfo {
	resp := httpGet("/api/rate?tokenType=all")
	jsonObj := new(jsonResp)
	err := json.Unmarshal(resp, &jsonObj)
	checkErr(err)
	if jsonObj.Code != 0 {
		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	rResp := new(rateResp)
	err = json.Unmarshal(jsonObj.Data, &rResp)
	checkErr(err)
	res := make([]*db.CryptocurrencyInfo, 0, 16)
	for _, coin := range rResp.AllCoin {
		res = append(res, &db.CryptocurrencyInfo{coin.Name, coin.Code, coin.Rate})
	}
	return res
}

type DepositResp struct {
	GoOn    bool               `json:"goon"`
	NextSeq int64              `json:"nextseq"`
	Deposit []db.DepositRecord `json:"deposit"`
}

func Deposite(req int64) *DepositResp {
	resp := httpGet(fmt.Sprintf("/api/deposite?req=%d", req))
	jsonObj := new(jsonResp)
	err := json.Unmarshal(resp, &jsonObj)
	checkErr(err)
	if jsonObj.Code != 0 {
		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	res := new(DepositResp)
	err = json.Unmarshal(jsonObj.Data, &res)
	checkErr(err)
	return res
}

type StatusResp struct {
	Seq       int64  `json:"seq"`
	UpdateAt  uint64 `json:"update_at"`
	Address   string `json:"address"`
	TokenType string `json:"tokenType"`
	StatusStr string `json:"status"`
	Status    int
}

const (
	str_waiting_deposit = "waiting_deposit"
	str_waiting_send    = "waiting_send"
	str_waiting_confirm = "waiting_confirm"
	str_done            = "done"
)

const (
	StatusWaitingDeposit = iota
	StatusWaitingSend
	StatusWaitingConfirm
	StatusDone
)

func Status(address string, currencyType string) []*StatusResp {
	resp := httpGet(fmt.Sprintf("/api/status?address=%s&coin_type=%s", address, currencyType))
	jsonObj := new(jsonResp)
	err := json.Unmarshal(resp, &jsonObj)
	checkErr(err)
	if jsonObj.Code != 0 {
		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	res := make([]*StatusResp, 0, 16)
	err = json.Unmarshal(jsonObj.Data, &res)
	checkErr(err)
	for _, s := range res {
		switch s.StatusStr {
		case str_waiting_deposit:
			s.Status = StatusWaitingDeposit
		case str_waiting_send:
			s.Status = StatusWaitingSend
		case str_waiting_confirm:
			s.Status = StatusWaitingConfirm
		case str_done:
			s.Status = StatusDone
		default:
			panic("wrong status string")
		}
	}
	return res
}

type SendCoinInfo struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
	Id      uint64 `json:"id"`
}

var logger *log.Logger

func init() {
	os.MkdirAll(config.GetDaemonConfig().LogFolder, 0755)
	f, err := os.OpenFile(config.GetDaemonConfig().LogFolder+"send-coin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	logger = log.New(f, "INFO", log.Ldate|log.Ltime)
}

func SendCoin(addrAndAmount []*SendCoinInfo) {
	body, err := json.Marshal(addrAndAmount)
	logger.Println(body)
	checkErr(err)
	resp := httpPost("/api/send-coin", body)
	jsonObj := new(jsonResp)
	err = json.Unmarshal(resp, &jsonObj)
	checkErr(err)
	if jsonObj.Code != 0 {
		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
}
