package main

import (
	"encoding/json"
	"fmt"
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service"
	"github.com/spaco/affiliate/src/service/db"
	"github.com/spaco/affiliate/src/tracking_code"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
)

func init() {
	os.MkdirAll(config.GetDaemonConfig().LogFolder, 0755)
	f, err := os.OpenFile(config.GetDaemonConfig().LogFolder+"task.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkErr(err)
	logger = log.New(f, "INFO", log.Ldate|log.Ltime)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var logger *log.Logger

func init() {
	os.MkdirAll(config.GetDaemonConfig().LogFolder, 0755)
	f, err := os.OpenFile(config.GetDaemonConfig().LogFolder+"server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	refCode := r.PostFormValue("ref")
	// fmt.Printf("Addr: %s, Ref: %s", addr, ref)
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
	renderCodeTemplate(w, "my-invitation", &struct{ Code string }{Code: "9527"})
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
		AllCurrency []*db.CryptocurrencyInfo
		Ref         string
	}{config.GetServerConfig().CoinName, service.AllCryptocurrency(), r.FormValue("ref")})
}

func checkStatusHandler(w http.ResponseWriter, r *http.Request) {
	transferAddress := "test" // TODO
	json.NewEncoder(w).Encode(&struct {
		TransferAddress string `json:"transferAddress"`
	}{transferAddress})
}

func getAddrHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	depositAddr := service.MappingDepositAddr(r.FormValue("address"), r.FormValue("currencyType"), r.FormValue("ref"))
	data := &struct {
		DepositAddr string `json:"depositAddr"`
	}{depositAddr}
	json.NewEncoder(w).Encode(&JsonObj{0, "", data})
}
