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

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service"
	"github.com/spaco/affiliate/src/service/db"
	client "github.com/spaco/affiliate/src/teller_client"
	"github.com/spaco/affiliate/src/tracking_code"
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
	fmt.Printf("Listening on :%d", config.Server.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil)
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
	// total := 0
	// if len(records)>0{
	// 	for _,r := range records{
	// 		total += r.SentAmount
	// 	}
	// }
	data := &struct {
		RewardRecords []db.RewardRecord `json:"records"`
		RewardRemain  uint64            `json:"remain"`
	}{records, service.QueryRewardRemain(addr)}
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

func buyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	renderBuyTemplate(w, "index", struct {
		CoinName    string
		AllCurrency []db.CryptocurrencyInfo
		Ref         string
	}{config.GetServerConfig().CoinName, service.AllCryptocurrency(), r.FormValue("ref")})
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
	res := service.QueryDepositRecord(addr)
	if len(res) > 0 {
		var totalDeposit, totalBuy uint64
		for _, dr := range res {
			totalDeposit += dr.DepositAmount
			totalBuy += dr.BuyAmount
		}
		data = fmt.Sprintf("found %d deposit, Total amount is %d, buy %d "+config.GetServerConfig().CoinName, len(res), totalDeposit, totalBuy)
	} else if service.CheckMappingAddr(addr, currencyType) {
		status, err := client.Status(addr, currencyType)
		if err != nil {
			json.NewEncoder(w).Encode(&JsonObj{2, "Teller api error: " + err.Error(), nil})
			return
		}
		if len(status) > 0 {
			data = fmt.Sprintf("found %d deposit, please waiting for confirm", len(status))
		}
		data = "Not found deposit record."
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
	if len(currencyType) == 0 || !service.CheckCryptocurrency(currencyType) {
		json.NewEncoder(w).Encode(&JsonObj{2, "Cryptocurrency type is not valid: " + currencyType, nil})
		return
	}
	depositAddr, err := service.MappingDepositAddr(addr, currencyType, r.PostFormValue("ref"))
	if err != nil {
		json.NewEncoder(w).Encode(&JsonObj{2, "Teller api error: " + err.Error(), nil})
		return
	}
	data := &struct {
		DepositAddr string `json:"depositAddr"`
	}{depositAddr}
	json.NewEncoder(w).Encode(&JsonObj{0, "", data})
}
