package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/db/postgresql"
	"github.com/spaco/affiliate/src/tracking_code"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic Error: %s", err)
		}
	}()
	http.HandleFunc("/code/", codeHandler)
	http.HandleFunc("/code/generate/", generateHandler)
	http.HandleFunc("/code/my-invitation/", myInvitationHandler)
	fsh := http.FileServer(http.Dir("s"))
	http.Handle("/s/", http.StripPrefix("/s/", fsh))
	http.HandleFunc("/favicon.ico", serveFileHandler)
	http.HandleFunc("/robots.txt", serveFileHandler)
	config := config.GetConfig()
	postgresql.OpenDb(&config.Db)
	defer postgresql.CloseDb()
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
	renderTemplate(w, "code", struct{ Ref string }{Ref: r.FormValue("ref")})
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addr := r.PostFormValue("address")
	ref := r.PostFormValue("ref")
	// fmt.Printf("Addr: %s, Ref: %s", addr, ref)
	id := postgresql.GetTrackingCodeOrGenerate(addr, ref)
	code := tracking_code.GenerateCode(id)
	server := config.GetConfig().Server
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
	renderTemplate(w, "generate", &struct {
		BuyUrl  string
		JoinUrl string
	}{contextPath + "/?ref=" + code, contextPath + "/code/?ref=" + code})
}

func myInvitationHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "my-invitation", &struct{ Code string }{Code: "9527"})
}

var templates = template.Must(template.ParseGlob("tpl/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
