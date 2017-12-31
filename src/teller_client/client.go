package teller_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	//	"math/rand"
	"crypto/md5"
	"github.com/shopspring/decimal"
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
	"net/http"
	"os"
	"time"
)

const print_req_resp = false

type jsonResp struct {
	Code   int16           `json:"code"`
	ErrMsg string          `json:"errmsg"`
	Data   json.RawMessage `json:"data"`
}

//func init() {
//	rand.Seed(time.Now().UnixNano())
//}
//
//var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
//
//func randStringRunes(n int) string {
//	b := make([]rune, n)
//	for i := range b {
//		b[i] = letterRunes[rand.Intn(len(letterRunes))]
//	}
//	return string(b)
//}
func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
func checkErr(err error, errPanic bool) error {
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return err
		}
	}
	return nil
}

func setAuthHeaders(req *http.Request, teller *config.Teller) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	h := md5.New()
	io.WriteString(h, timestamp+teller.ApiToken)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("auth", fmt.Sprintf("%x", h.Sum(nil)))
}

func httpReq(errPanic bool, url string, method string, reqBody io.Reader, contentType string) ([]byte, error) {
	conf := config.GetServerConfig()
	client := &http.Client{}
	req, err := http.NewRequest(method, conf.Teller.ContextPath+url, reqBody)
	if resErr := checkErr(err, errPanic); resErr != nil {
		return nil, resErr
	}
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}
	setAuthHeaders(req, &conf.Teller)
	resp, err := client.Do(req)
	if resErr := checkErr(err, errPanic); resErr != nil {
		return nil, resErr
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resErr := checkErr(err, errPanic); resErr != nil {
		return nil, resErr
	}
	if print_req_resp {
		fmt.Printf(string(body))
	}
	return body, nil
}
func httpGet(errPanic bool, url string) ([]byte, error) {
	return httpReq(errPanic, url, "GET", nil, "")
}
func httpPost(errPanic bool, url string, json []byte) ([]byte, error) {
	return httpReq(errPanic, url, "POST", bytes.NewBuffer(json), "application/json")
}

type bindResp struct {
	Address      string `json:"tokenAddress"`
	CurrencyType string `json:"tokenType"`
}

func Bind(currencyType string, address string) (string, error) {
	resp, err := httpPost(false, "/api/bind", []byte(fmt.Sprintf(`{"address":"%s","tokenType":"%s"}`, address, currencyType)))
	if err != nil {
		return "", err
	}
	jsonObj := new(jsonResp)
	err = json.Unmarshal(resp, &jsonObj)
	if err != nil {
		return "", err
	}
	if jsonObj.Code != 0 {
		return "", errors.New(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	bResp := new(bindResp)
	err = json.Unmarshal(jsonObj.Data, &bResp)
	if err != nil {
		return "", err
	}
	if bResp.CurrencyType != currencyType {
		return "", errors.New("return not same currency type")
	}
	return bResp.Address, nil
	//	return randStringRunes(32)
}

type rateResp struct {
	TokenType string     `json:"tokenType"`
	Rate      float32    `json:"rate"`
	AllCoin   []coinResp `json:"allcoin"`
}

type coinResp struct {
	Name      string  `json:"coin_name"`
	Code      string  `json:"coin_code"`
	Rate      float64 `json:"coin_rate"`
	UnitPower int32   `json:"unit"`
}

func Rate() []db.CryptocurrencyInfo {
	resp, _ := httpGet(true, "/api/rate?tokenType=all")
	jsonObj := new(jsonResp)
	err := json.Unmarshal(resp, &jsonObj)
	panicErr(err)
	if jsonObj.Code != 0 {
		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	rResp := new(rateResp)
	err = json.Unmarshal(jsonObj.Data, &rResp)
	panicErr(err)
	res := make([]db.CryptocurrencyInfo, 0, 16)
	for _, coin := range rResp.AllCoin {
		res = append(res, db.CryptocurrencyInfo{coin.Name, coin.Code, decimal.NewFromFloat(coin.Rate).String(), coin.UnitPower})
	}
	return res
}

type DepositResp struct {
	GoOn     bool               `json:"goon"`
	NextSeq  int64              `json:"nextseq"`
	Deposits []db.DepositRecord `json:"deposits"`
}

func Deposite(req int64) *DepositResp {
	resp, _ := httpGet(true, fmt.Sprintf("/api/deposits?seq=%d", req))
	//	resp, _ := httpPost(true, "/api/deposits", []byte(fmt.Sprintf(`{"req":"%d"}`, req)))
	jsonObj := new(jsonResp)
	err := json.Unmarshal(resp, &jsonObj)
	panicErr(err)
	if jsonObj.Code != 0 {
		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	res := new(DepositResp)
	err = json.Unmarshal(jsonObj.Data, &res)
	panicErr(err)
	if len(res.Deposits) > 0 {
		for _, dr := range res.Deposits {
			dr.Rate = decimal.NewFromFloat(dr.RateFloat).String()
		}
	}
	return res
}

type statusesResp struct {
	Statuses []StatusResp `json:"statuses"`
}
type StatusResp struct {
	Seq       int64  `json:"seq"`
	UpdateAt  uint64 `json:"update_at"`
	Address   string `json:"address"`
	TokenType string `json:"tokenType"`
	StatusStr string `json:"status"`
	Status    int    `json:"-"`
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

func Status(address string, currencyType string) ([]StatusResp, error) {
	resp, err := httpGet(false, fmt.Sprintf("/api/status?address=%s&tokenType=%s", address, currencyType))
	if err != nil {
		return nil, err
	}
	jsonObj := new(jsonResp)
	err = json.Unmarshal(resp, &jsonObj)
	if err != nil {
		return nil, err
	}
	if jsonObj.Code != 0 {
		return nil, errors.New(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	res := new(statusesResp)
	err = json.Unmarshal(jsonObj.Data, &res)
	panicErr(err)
	for _, s := range res.Statuses {
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
			return nil, errors.New("wrong status string")
		}
	}
	return res.Statuses, nil
}

var sendCoinLogger *log.Logger

func getSendCoinLogger() *log.Logger {
	if sendCoinLogger == nil {
		os.MkdirAll(config.GetDaemonConfig().LogFolder, 0755)
		f, err := os.OpenFile(config.GetDaemonConfig().LogFolder+"send-coin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		panicErr(err)
		sendCoinLogger = log.New(f, "INFO", log.Ldate|log.Ltime)
	}
	return sendCoinLogger
}

type sendCoinInfo struct {
	Id         uint64 `json:"id"`
	Address    string `json:"address"`
	SentAmount uint64 `json:"amount"`
}

func SendCoin(addrAndAmount []db.RewardRecord) {
	arr := make([]sendCoinInfo, 0, len(addrAndAmount))
	for _, rr := range addrAndAmount {
		arr = append(arr, sendCoinInfo{rr.Id, rr.Address, rr.SentAmount})
	}
	body, err := json.Marshal(arr)
	if print_req_resp {
		fmt.Printf(string(body))
	}
	getSendCoinLogger().Println(body)
	panicErr(err)
	resp, _ := httpPost(true, "/api/send-coin", body)
	jsonObj := new(jsonResp)
	err = json.Unmarshal(resp, &jsonObj)
	panicErr(err)
	if jsonObj.Code != 0 {
		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
}
