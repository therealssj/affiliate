package teller_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	//	"math/rand"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	//	"github.com/shopspring/decimal"
	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
)

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
	//	timestamp := strconv.Itoa(time.Now().Unix())
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	//	hash := md5.New()
	//	io.WriteString(hash, timestamp+teller.ApiToken)
	hash := hmac.New(sha256.New, []byte(teller.ApiToken))
	hash.Write([]byte(timestamp))

	req.Header.Set("timestamp", timestamp)
	//	req.Header.Set("auth", fmt.Sprintf("%x", hash.Sum(nil)))
	req.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))
	req.Header.Set("affiliate", "true")
}

func httpReq(tellerConf *config.Teller, errPanic bool, url string, method string, reqBody io.Reader, contentType string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, tellerConf.ContextPath+url, reqBody)
	if resErr := checkErr(err, errPanic); resErr != nil {
		return nil, resErr
	}
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}
	setAuthHeaders(req, tellerConf)
	resp, err := client.Do(req)
	if resErr := checkErr(err, errPanic); resErr != nil {
		return nil, resErr
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resErr := checkErr(err, errPanic); resErr != nil {
		return nil, resErr
	}
	if tellerConf.Debug {
		fmt.Printf(string(body))
	}
	return body, nil
}
func httpGet(tellerConf *config.Teller, errPanic bool, url string) ([]byte, error) {
	return httpReq(tellerConf, errPanic, url, "GET", nil, "")
}
func httpPost(tellerConf *config.Teller, errPanic bool, url string, json []byte) ([]byte, error) {
	return httpReq(tellerConf, errPanic, url, "POST", bytes.NewBuffer(json), "application/json")
}

type bindResp struct {
	Address      string `json:"tokenAddress"`
	CurrencyType string `json:"tokenType"`
}

func Bind(currencyType string, address string) (string, error) {
	conf := config.GetServerConfig()
	if conf.TestMode {
		return address + "-" + currencyType + "-bind-mock-result", nil
	}
	resp, err := httpPost(&(conf.Teller), false, "/api/bind", []byte(fmt.Sprintf(`{"address":"%s","tokenType":"%s"}`, address, currencyType)))
	if err != nil {
		return "", err
	}
	return bindRespProcess(resp, currencyType)
}

func bindRespProcess(response []byte, currencyType string) (string, error) {
	jsonObj := new(jsonResp)
	err := json.Unmarshal(response, &jsonObj)
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
}

//type rateResp struct {
//	TokenType string     `json:"tokenType"`
//	Rate      float32    `json:"rate"`
//	AllCoin   []coinResp `json:"allcoin"`
//}
//
//type coinResp struct {
//	Name      string  `json:"coin_name"`
//	Code      string  `json:"coin_code"`
//	Rate      float64 `json:"coin_rate"`
//	UnitPower int32   `json:"unit"`
//}

//func Rate() []db.CryptocurrencyInfo {
//	resp, _ := httpGet(&(config.GetDaemonConfig().Teller), true, "/api/rate?tokenType=all")
//	jsonObj := new(jsonResp)
//	err := json.Unmarshal(resp, &jsonObj)
//	panicErr(err)
//	if jsonObj.Code != 0 {
//		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
//	}
//	rResp := new(rateResp)
//	err = json.Unmarshal(jsonObj.Data, &rResp)
//	panicErr(err)
//	res := make([]db.CryptocurrencyInfo, 0, len(rResp.AllCoin))
//	for _, coin := range rResp.AllCoin {
//		res = append(res, db.CryptocurrencyInfo{coin.Name, coin.Code, decimal.NewFromFloat(coin.Rate).String(), coin.UnitPower})
//	}
//	return res
//}

//func RateWithErr() ([]db.CryptocurrencyInfo, error) {
//	resp, err := httpGet(&(config.GetServerConfig().Teller), false, "/api/rate?tokenType=all")
//	if err != nil {
//		return nil, err
//	}
//	jsonObj := new(jsonResp)
//	err = json.Unmarshal(resp, &jsonObj)
//	if err != nil {
//		return nil, err
//	} else if jsonObj.Code != 0 {
//		return nil, errors.New((fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code)))
//	}
//	rResp := new(rateResp)
//	err = json.Unmarshal(jsonObj.Data, &rResp)
//	if err != nil {
//		return nil, err
//	}
//	slice := make([]db.CryptocurrencyInfo, 0, len(rResp.AllCoin))
//	for _, coin := range rResp.AllCoin {
//		slice = append(slice, db.CryptocurrencyInfo{coin.Name, coin.Code, decimal.NewFromFloat(coin.Rate).String(), coin.UnitPower})
//	}
//	return slice, nil
//}

//type DepositResp struct {
//	GoOn     bool               `json:"goon"`
//	NextSeq  int64              `json:"nextseq"`
//	Deposits []db.DepositRecord `json:"deposits"`
//}
//
//func Deposite(req int64) *DepositResp {
//	tellerConf := config.GetDaemonConfig().Teller
//	if tellerConf.Debug {
//		fmt.Printf("seq:%d", req)
//	}
//	resp, _ := httpGet(&tellerConf, true, fmt.Sprintf("/api/deposits?seq=%d", req))
//	//	resp, _ := httpPost(true, "/api/deposits", []byte(fmt.Sprintf(`{"req":"%d"}`, req)))
//	jsonObj := new(jsonResp)
//	err := json.Unmarshal(resp, &jsonObj)
//	panicErr(err)
//	if jsonObj.Code != 0 {
//		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
//	}
//	res := new(DepositResp)
//	err = json.Unmarshal(jsonObj.Data, &res)
//	panicErr(err)
//	if len(res.Deposits) > 0 {
//		for i, _ := range res.Deposits {
//			res.Deposits[i].Rate = decimal.NewFromFloat(res.Deposits[i].RateFloat).String()
//		}
//	}
//	return res
//}

type statusesResp struct {
	Statuses []StatusResp `json:"statuses"`
}
type StatusResp struct {
	Seq       int64  `json:"seq"`
	UpdateAt  uint64 `json:"updated_at"`
	TokenType string `json:"coin_type"`
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
	resp, err := httpGet(&(config.GetServerConfig().Teller), false, fmt.Sprintf("/api/status?address=%s&tokenType=%s", address, currencyType))
	if err != nil {
		return nil, err
	}
	return statusRespProcess(resp, currencyType)

}

//func statusRespTidy(arr []StatusResp) string {
//	var waitingDeposit, waitingSend, waitingConfirm, done uint32
//	for _, s := range arr {
//		switch s.Status {
//		case StatusWaitingDeposit:
//			waitingDeposit++
//		case StatusWaitingSend:
//			waitingSend++
//		case StatusWaitingConfirm:
//			waitingConfirm++
//		case StatusDone:
//			done++
//		}
//	}
//	if done =0{
//		return "";
//	}
//	return "";
//}

func statusRespProcess(response []byte, currencyType string) ([]StatusResp, error) {
	jsonObj := new(jsonResp)
	err := json.Unmarshal(response, &jsonObj)
	if err != nil {
		return nil, err
	}
	if jsonObj.Code != 0 {
		return nil, errors.New(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	resp := new(statusesResp)
	err = json.Unmarshal(jsonObj.Data, &resp)
	panicErr(err)
	res := make([]StatusResp, 0, len(resp.Statuses))
	for _, s := range resp.Statuses {
		if s.TokenType != currencyType {
			continue
		}
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
			return nil, errors.New("wrong status: " + s.StatusStr)
		}
		res = append(res, s)
	}
	return res, nil
}

//var sendCoinLogger *log.Logger
//
//func getSendCoinLogger() *log.Logger {
//	if sendCoinLogger == nil {
//		logFolder := config.GetDaemonConfig().LogFolder
//		os.MkdirAll(logFolder, 0755)
//		f, err := os.OpenFile(logFolder+"send-coin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//		panicErr(err)
//		sendCoinLogger = log.New(f, "INFO", log.Ldate|log.Ltime)
//	}
//	return sendCoinLogger
//}
//
//type sendCoinInfo struct {
//	Id         uint64 `json:"id"`
//	Address    string `json:"address"`
//	SentAmount uint64 `json:"amount"`
//}
//
//func SendCoin(addrAndAmount []db.RewardRecord) {
//	slice := make([]sendCoinInfo, 0, len(addrAndAmount))
//	for _, rr := range addrAndAmount {
//		slice = append(slice, sendCoinInfo{rr.Id, rr.Address, rr.SentAmount})
//	}
//	body, err := json.Marshal(slice)
//	tellerConf := config.GetDaemonConfig().Teller
//	if tellerConf.Debug {
//		fmt.Printf(string(body))
//	}
//	getSendCoinLogger().Println(string(body))
//	panicErr(err)
//	resp, _ := httpPost(&tellerConf, true, "/api/send-coin", body)
//	jsonObj := new(jsonResp)
//	err = json.Unmarshal(resp, &jsonObj)
//	panicErr(err)
//	if jsonObj.Code != 0 {
//		panic(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
//	}
//}
type configResp struct {
	Enabled       bool                `json:"enabled"`
	MaxBoundAddrs uint32              `json:"max_bound_addrs"`
	MaxDecimals   uint32              `json:"max_decimals"`
	AllCoins      map[string]coinInfo `json:"allcoins"`
}

type coinInfo struct {
	Name                  string `json:"coin_name"`
	Rate                  string `json:"rate"`
	Unit                  uint64 `json:"unit"`
	Enabled               bool   `json:"enabled"`
	ConfirmationsRequired uint32 `json:"confirmations_required"`
}

func Config() (*configResp, error) {
	resp, err := httpGet(&(config.GetServerConfig().Teller), false, "/api/config")
	if err != nil {
		return nil, err
	}
	return configRespProcess(resp)
}

func configRespProcess(response []byte) (*configResp, error) {
	jsonObj := new(jsonResp)
	err := json.Unmarshal(response, &jsonObj)
	if err != nil {
		return nil, err
	}
	if jsonObj.Code != 0 {
		return nil, errors.New(fmt.Sprintf("%s code:%d", jsonObj.ErrMsg, jsonObj.Code))
	}
	confResp := new(configResp)
	err = json.Unmarshal(jsonObj.Data, &confResp)
	if err != nil {
		return nil, err
	}
	return confResp, nil
}

func RateWithErr() ([]db.CryptocurrencyInfo, error) {
	conf := config.GetServerConfig()
	if conf.TestMode {
		json := `{
			"errmsg": "",
			"code": 0,
			"Data": {
				"enabled": true,
				"max_bound_addrs": 4,
				"max_decimals": 3,
				"allcoins": {
					"BTC": {
						"coin_name": "BTC",
						"rate": "55556",
						"enabled": true,
						"unit": 100000000,
						"confirmations_required": 1
					},
					"ETH": {
						"coin_name": "ETH",
						"rate": "5912",
						"enabled": true,
						"unit": 1000000000,
						"confirmations_required": 4
					},
					"SKY": {
						"coin_name": "SKY",
						"rate": "126",
						"enabled": true,
						"unit": 1000000,
						"confirmations_required": 0
					},
					"XMR": {
						"coin_name": "XMR",
						"rate": "1601",
						"enabled": true,
						"unit": 1000000000000,
						"confirmations_required": 4
					}
				}
			}
		}`
		confResp, _ := configRespProcess([]byte(json))
		return rateWithErrProcess(confResp)
	}
	confResp, err := Config()
	if err != nil {
		return nil, err
	}
	return rateWithErrProcess(confResp)
}

func rateWithErrProcess(confResp *configResp) ([]db.CryptocurrencyInfo, error) {
	slice := make([]db.CryptocurrencyInfo, 0, len(confResp.AllCoins))
	for _, coin := range confResp.AllCoins {
		slice = append(slice, db.CryptocurrencyInfo{coin.Name, coin.Name, coin.Rate, int32(math.Log10(float64(coin.Unit))), coin.Enabled})
	}
	return slice, nil
}

type StatsLeftInfo struct {
	TotalHours   string  `json:"-"` //`json:"total_hours"`
	SoldRatio    float64 `json:"sold_ratio"`
	RewardHours  string  `json:"-"` //`json:"reward_hours"`
	TotalAmount  string  `json:"total_amount"`
	RewardAmount string  `json:"reward_amount"`
}

func StatsLeft() (*StatsLeftInfo, error) {
	conf := config.GetServerConfig()
	if conf.TestMode {
		return statsLeftRespProcess([]byte(`{"total_hours": "18751304", "sold_ratio": 0.98, "reward_hours": "18751304", "total_amount": "110357.663000", "reward_amount": "110357.663000"}`))
	} else {
		json, err := httpGet(&(config.GetServerConfig().Teller), false, "/stats/left")
		if err != nil {
			return nil, err
		}
		return statsLeftRespProcess(json)
	}
}

func statsLeftRespProcess(response []byte) (*StatsLeftInfo, error) {
	statsLeftInfo := new(StatsLeftInfo)
	err := json.Unmarshal(response, &statsLeftInfo)
	if err != nil {
		return nil, err
	}
	return statsLeftInfo, nil
}
