package main

import (
	"fmt"
	"github.com/koding/multiconfig"
	"github.com/spaco/affiliate/src/tracking_code"
	"html/template"
	"net/http"
	"path"
)

type Config struct {
	DbHost     string `default:"localhost"`
	DbPort     string `default:"5432"`
	DbUser     string `default:"lijt"`
	DbPassword string `default:"lijtlijt"`
	DbName     string `default:"affiliate"`
	DbSslMode  string `default:"disable"`
	ListenPort string `default:":6060"`
}

func getConfig() *Config {
	m := multiconfig.NewWithPath("config.toml") // supports TOML, JSON and YAML
	// Get an empty struct for your configuration
	serverConf := new(Config)
	// Populated the serverConf struct
	err := m.Load(serverConf) // Check for error
	if err != nil {
		fmt.Println("ERROR: %s", err)
	}
	m.MustLoad(serverConf) // Panic's if there is any error
	return serverConf
}
func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR: %s", err)
		}
	}()
	http.HandleFunc("/code/", codeHandler)
	http.HandleFunc("/code/generate/", generateHandler)
	http.HandleFunc("/code/my-invitation/", myInvitationHandler)
	fsh := http.FileServer(http.Dir("s"))
	http.Handle("/s/", http.StripPrefix("/s/", fsh))
	http.HandleFunc("/favicon.ico", serveFileHandler)
	http.HandleFunc("/robots.txt", serveFileHandler)
	config := getConfig()
	fmt.Println("Listening on", config.ListenPort)
	err := http.ListenAndServe(config.ListenPort, nil)
	if err != nil {
		println("error： %s", err)
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
	fmt.Printf("Addr: %s, Ref: %s", addr, ref)
	code := tracking_code.GenerateCode(12)
	renderTemplate(w, "generate", &struct{ Code string }{Code: code})
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
