package main

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service"
	"github.com/spaco/affiliate/src/service/db"
	client "github.com/spaco/affiliate/src/teller_client"
	"github.com/spaco/affiliate/src/tracking_code"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"sort"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var logger *log.Logger

func init() {
	os.MkdirAll(config.GetServerConfig().LogFolder, 0755)
	f, err := os.OpenFile(config.GetServerConfig().LogFolder+"server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	logger = log.New(f, "INFO", log.Ldate|log.Ltime)
}

func convertUnitPower(val uint64, unitPower int32) string {
	return decimal.New(int64(val), 0).Div(decimal.New(1, unitPower)).String()
}

func convertUnit(val uint64) string {
	return convertUnitPower(val, config.BUY_COIN_UNIT_POWER)
}

func convertOfCurrency(val uint64, currencyType string) string {
	var unitPower int32
	if info, ok := getAllCryptocurrencyMap()[currencyType]; ok {
		unitPower = info.UnitPower
	} else {
		info := service.GetCryptocurrency(currencyType)
		if info == nil {
			panic(currencyType + " is not support.")
		}
		unitPower = info.UnitPower
	}
	return convertUnitPower(val, unitPower)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic Error: %s", err)
			debug.PrintStack()
			logger.Println(debug.Stack())
		}
	}()
	http.HandleFunc("/", buyHandler)
	http.HandleFunc("/get-address/", getAddrHandler)
	http.HandleFunc("/check-status/", checkStatusHandler)
	http.HandleFunc("/code/", codeHandler)
	http.HandleFunc("/code/generate/", generateHandler)
	http.HandleFunc("/code/my-invitation/", myInvitationHandler)
	fsh := http.FileServer(http.Dir("s"))
	http.Handle("/s/", http.StripPrefix("/s/", fsh))
	http.HandleFunc("/favicon.ico", serveFileHandler)
	http.HandleFunc("/robots.txt", serveFileHandler)
	config := config.GetServerConfig()
	db.OpenDb(&config.Db)
	defer db.CloseDb()
	fmt.Printf("Listening on :%d", config.Server.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Server.ListenPort), nil)
	if err != nil {
		println("ListenAndServe Error： %s", err)
	}
}
func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	fname := path.Base(r.URL.Path)
	http.ServeFile(w, r, "./s/"+fname)
}
func codeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	renderCodeTemplate(w, "index", struct{ Ref string }{Ref: r.FormValue("ref")})
}

type JsonObj struct {
	Code   uint8       `json:"code"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr := r.PostFormValue("address")
	if _, err := cipher.DecodeBase58Address(addr); err != nil {
		json.NewEncoder(w).Encode(&JsonObj{1, addr + " is not valid. " + err.Error(), nil})
		return
	}
	refCode := r.PostFormValue("ref")
	id := service.GetTrackingCodeOrGenerate(addr, refCode)
	code := tracking_code.GenerateCode(id)
	server := config.GetServerConfig().Server
	contextPath := "http"
	if server.Https {
		contextPath = "https"
	}
	contextPath += "://" + server.Domain
	if server.Https {
		if server.Port != 443 {
			contextPath = fmt.Sprintf("%s:%d", contextPath, server.Port)
		}
	} else {
		if server.Port != 80 {
			contextPath = fmt.Sprintf("%s:%d", contextPath, server.Port)
		}
	}
	data := &struct {
		BuyUrl  string `json:"buyUrl"`
		JoinUrl string `json:"joinUrl"`
	}{contextPath + "/?ref=" + code, contextPath + "/code/?ref=" + code}
	json.NewEncoder(w).Encode(&JsonObj{0, "", data})
}

func myInvitationHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr := r.PostFormValue("address")
	if _, err := cipher.DecodeBase58Address(addr); err != nil {
		json.NewEncoder(w).Encode(&JsonObj{1, addr + " is not valid. " + err.Error(), nil})
		return
	}
	records := service.QueryRewardRecord(addr)
	if len(records) > 0 {
		for i, _ := range records {
			records[i].CalAmountStr = convertUnit(records[i].CalAmount)
			records[i].SentAmountStr = convertUnit(records[i].SentAmount)
		}
	}
	data := &struct {
		CoinName      string            `json:"coinName"`
		RewardRecords []db.RewardRecord `json:"records"`
		RewardRemain  string            `json:"remain"`
	}{config.GetServerConfig().CoinName, records, convertUnit(service.QueryRewardRemain(addr))}
	json.NewEncoder(w).Encode(&JsonObj{0, "", data})
}

var codeTemplates = template.Must(template.ParseGlob("tpl-code/*.html"))

func renderCodeTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := codeTemplates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

var buyTemplates = template.Must(template.ParseGlob("tpl-buy/*.html"))

func renderBuyTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := buyTemplates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

var allCryptocurrencyMap map[string]db.CryptocurrencyInfo
var allCurrencyMapLastUpdated int64

func getAllCryptocurrencyMap() map[string]db.CryptocurrencyInfo {
	if allCryptocurrencyMap == nil || len(allCryptocurrencyMap) == 0 || time.Now().Unix()-allCurrencyMapLastUpdated > 3600 {
		allCryptocurrencyMap = service.AllCryptocurrencyMap()
		allCurrencyMapLastUpdated = time.Now().Unix()
		return allCryptocurrencyMap
	}
	return allCryptocurrencyMap
}

func allCryptocurrency() []db.CryptocurrencyInfo {
	m := getAllCryptocurrencyMap()
	res := make([]db.CryptocurrencyInfo, 0, len(m))
	for _, value := range m {
		res = append(res, value)
	}
	sort.Sort(db.CryptocurrencyInfoSlice(res))
	return res
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	renderBuyTemplate(w, "index", struct {
		CoinName    string
		AllCurrency []db.CryptocurrencyInfo
		Ref         string
	}{config.GetServerConfig().CoinName, allCryptocurrency(), r.FormValue("ref")})
}

func checkStatusHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr := r.PostFormValue("address")
	if _, err := cipher.DecodeBase58Address(addr); err != nil {
		json.NewEncoder(w).Encode(&JsonObj{1, addr + " is not valid. " + err.Error(), nil})
		return
	}
	currencyType := r.PostFormValue("currencyType")
	var data string
	res := service.QueryDepositRecord(addr, currencyType)
	if len(res) > 0 {
		var totalDeposit, totalBuy uint64
		for _, dr := range res {
			totalDeposit += dr.DepositAmount
			totalBuy += dr.BuyAmount
		}
		data = fmt.Sprintf("found %d deposit, Total amount is %s, buy %s %s", len(res),
			convertOfCurrency(totalDeposit, currencyType), convertUnit(totalBuy), config.GetServerConfig().CoinName)
	} else if service.CheckMappingAddr(addr, currencyType) {
		status, err := client.Status(addr, currencyType)
		if err != nil {
			json.NewEncoder(w).Encode(&JsonObj{2, "Teller api error: " + err.Error(), nil})
			return
		}
		if len(status) > 0 {
			data = fmt.Sprintf("found %d deposit, please waiting for confirm", len(status))
		} else {
			data = "Not found deposit record."
		}
	} else {
		data = "Not found deposit record."
	}
	json.NewEncoder(w).Encode(&JsonObj{0, "", data})
}

func getAddrHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr := r.PostFormValue("address")
	if _, err := cipher.DecodeBase58Address(addr); err != nil {
		json.NewEncoder(w).Encode(&JsonObj{1, addr + " is not valid. " + err.Error(), nil})
		return
	}
	currencyType := r.PostFormValue("currencyType")
	if _, ok := getAllCryptocurrencyMap()[currencyType]; !ok {
		json.NewEncoder(w).Encode(&JsonObj{2, "Cryptocurrency type is not valid: " + currencyType, nil})
		return
	}
	depositAddr, err := service.MappingDepositAddr(addr, currencyType, r.PostFormValue("ref"))
	if err != nil {
		json.NewEncoder(w).Encode(&JsonObj{2, "Teller api error: " + err.Error(), nil})
		return
	}
	json.NewEncoder(w).Encode(&JsonObj{0, "", &struct {
		DepositAddr string `json:"depositAddr"`
	}{depositAddr}})
}
