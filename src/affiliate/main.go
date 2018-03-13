package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"sort"
	"time"

	"github.com/robfig/cron"
	"github.com/shopspring/decimal"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service"
	"github.com/spolabs/affiliate/src/service/db"
	client "github.com/spolabs/affiliate/src/teller_client"
	"github.com/spolabs/affiliate/src/tracking_code"
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
	return convertUnitPower(val, int32(config.GetServerConfig().CoinUnitPower))
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
	// http.HandleFunc("/stats-left/", statsLefthandler)
	http.HandleFunc("/qr-code/", qrCodehandler)
	http.HandleFunc("/get-address/", getAddrHandler)
	http.HandleFunc("/check-status/", checkStatusHandler)
	http.HandleFunc("/get-rate/", getRateHandler)
	// http.HandleFunc("/code/", codeHandler)
	http.HandleFunc("/generate/", generateHandler)
	http.HandleFunc("/my-invitation/", myInvitationHandler)
	http.HandleFunc("/more-invitation/", moreInvitationHandler)
	fsh := http.FileServer(http.Dir("s"))
	http.Handle("/s/", cache(http.StripPrefix("/s/", fsh)))
	http.HandleFunc("/favicon.ico", serveFileHandler)
	http.HandleFunc("/robots.txt", serveFileHandler)
	config := config.GetServerConfig()
	db.OpenDb(&config.Db)
	defer db.CloseDb()
	c := cron.New()
	c.AddFunc("0 * * * * *", updateAllCryptocurrencySlice)
	c.AddFunc("0 * * * * *", refreshSoldRatio)
	c.AddFunc("10,20,30,40,50 40-42 11 * * *", updateAllCryptocurrencySlice)
	c.Start()
	fmt.Printf("Listening on %s:%d", config.Server.ListenIp, config.Server.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.ListenIp, config.Server.ListenPort), nil)
	if err != nil {
		println("ListenAndServe Error： %s", err)
	}
	c.Stop()
}
func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	fname := path.Base(r.URL.Path)
	w.Header().Set("Cache-Control", "max-age=604800") //7days
	http.ServeFile(w, r, "./s/"+fname)
}
func cache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=7776000") //90days
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// func codeHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	renderCodeTemplate(w, "index", struct {
// 		CoinName string
// 		Ref      string
// 	}{config.GetServerConfig().CoinName, r.FormValue("ref")})
// }

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
	id := service.GetTrackingCodeOrGenerate(addr, getRefCookie(r))
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
	times := len(records)
	if times > 4 {
		records = records[0:4]
	}
	if len(records) > 0 {
		for i, _ := range records {
			records[i].CalAmountStr = convertUnit(records[i].CalAmount)
			records[i].SentAmountStr = convertUnit(records[i].SentAmount)
		}
	}
	data := &struct {
		CoinName      string            `json:"coinName"`
		Address       string            `json:"address"`
		Times         int               `json:"times"`
		RewardRecords []db.RewardRecord `json:"records"`
		RewardRemain  string            `json:"remain"`
	}{config.GetServerConfig().CoinName, addr, times, records, convertUnit(service.QueryRewardRemain(addr))}
	json.NewEncoder(w).Encode(&JsonObj{0, "", data})
}

func moreInvitationHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr := r.FormValue("address")
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
	renderTemplate(w, "more_invitation", struct {
		CoinName      string
		Times         int
		RewardRecords []db.RewardRecord
		RewardRemain  string
	}{config.GetServerConfig().CoinName, len(records), records, convertUnit(service.QueryRewardRemain(addr))})
}

var codeTemplates = template.Must(template.ParseGlob("tpl-code/*.html"))

// func renderCodeTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
// 	err := codeTemplates.ExecuteTemplate(w, tmpl+".html", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

var templates = template.Must(template.ParseGlob("*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

var allCryptocurrencySlice []Cryptocurrency

func updateAllCryptocurrencySlice() {
	slice, err := client.RateWithErr()
	if err != nil {
		logger.Printf("client.RateWithErr() err: %s", err)
		slice = service.AllCryptocurrency()
	} else {
		service.SyncCryptocurrency(slice)
	}
	sort.Sort(db.CryptocurrencyInfoSlice(slice))
	sl := make([]Cryptocurrency, 0, len(slice))
	for _, info := range slice {
		sl = append(sl, newCryptocurrency(&info))
	}
	allCryptocurrencySlice = sl
}

func getAllCryptocurrencyMap() map[string]db.CryptocurrencyInfo {
	slice := allCryptocurrency()
	m := make(map[string]db.CryptocurrencyInfo, len(slice))
	for _, info := range slice {
		m[info.ShortName] = info.CryptocurrencyInfo
	}
	return m
}

type Cryptocurrency struct {
	db.CryptocurrencyInfo
	ReverseRate string `json:"reverse_rate"`
}

func newCryptocurrency(info *db.CryptocurrencyInfo) Cryptocurrency {
	de, _ := decimal.NewFromString(info.Rate)
	return Cryptocurrency{*info, decimal.NewFromFloat(1).DivRound(de, info.UnitPower).String()}
}

func allCryptocurrency() []Cryptocurrency {
	if allCryptocurrencySlice == nil || len(allCryptocurrencySlice) == 0 {
		updateAllCryptocurrencySlice()
	}
	return allCryptocurrencySlice
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ref := r.FormValue("ref")
	if len(ref) > 0 {
		setRefCookie(w, ref)
	}
	if !statsLeftInit {
		refreshSoldRatio()
		statsLeftInit = true
	}
	renderTemplate(w, "index", struct {
		CoinName    string
		AllCurrency []Cryptocurrency
		StatsLeft   client.StatsLeftInfo
	}{config.GetServerConfig().CoinName, allCryptocurrency(), *statsLeftInfo})
}

func checkStatusHandler(w http.ResponseWriter, r *http.Request) {
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
		//	} else if service.CheckMappingAddr(addr, currencyType) {
		//		status, err := client.Status(addr, currencyType)
		//		if err != nil {
		//			json.NewEncoder(w).Encode(&JsonObj{2, "Teller api error: " + err.Error(), nil})
		//			return
		//		}
		//		if len(status) > 0 {
		//			data = fmt.Sprintf("found %d deposit, please waiting for confirm", len(status))
		//		} else {
		//			data = "Not found deposit record."
		//		}
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
	depositAddr, first, err := service.MappingDepositAddr(addr, currencyType, getRefCookie(r))
	if err != nil {
		json.NewEncoder(w).Encode(&JsonObj{2, "Teller api error: " + err.Error(), nil})
		return
	}
	json.NewEncoder(w).Encode(&JsonObj{0, "", &struct {
		DepositAddr string `json:"depositAddr"`
		First       bool   `json:"first"`
	}{depositAddr, first}})
}

func getRateHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&JsonObj{0, "", allCryptocurrency()})
}

func qrCodehandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	content := r.FormValue("content")
	if len(content) == 0 {
		http.Error(w, "parameter content is blank", http.StatusInternalServerError)
		return
	}
	png, err := qrcode.Encode(content, qrcode.Medium, -5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=2592000") //30days
	w.Write(png)
}

var statsLeftInfo *client.StatsLeftInfo
var statsLeftInit = false

func refreshSoldRatio() {
	info, err := client.StatsLeft()
	if err != nil {
		logger.Printf("client.StatsLeft() err: %s", err)
	} else {
		statsLeftInfo = info
	}
}

func statsLefthandler(w http.ResponseWriter, r *http.Request) {
	if !statsLeftInit {
		refreshSoldRatio()
		statsLeftInit = true
	}
	json.NewEncoder(w).Encode(&JsonObj{0, "", statsLeftInfo})
}

const cookie_name = "ref"

func setRefCookie(w http.ResponseWriter, ref string) {
	// cookie will get expired after 1 year
	expires := time.Now().AddDate(1, 0, 0)
	cookie := http.Cookie{
		Name:    cookie_name,
		Value:   ref,
		Path:    "/",
		Expires: expires,
	}
	http.SetCookie(w, &cookie)
}

func getRefCookie(r *http.Request) string {
	var cookie, err = r.Cookie(cookie_name)
	if err == nil {
		return cookie.Value
	}
	return ""
}
