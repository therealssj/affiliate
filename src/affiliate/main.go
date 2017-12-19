package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/spaco/affiliate/src/tracking_code"

	rice "github.com/GeertJohan/go.rice"
	"github.com/julienschmidt/httprouter"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR: %s", err)
		}
	}()
	listen := flag.String("-listen", ":6060", "Interface and port to listen on")
	flag.Parse()
	fmt.Println("Listening on", *listen)
	log.Fatal(http.ListenAndServe(*listen, getRouter()))
}

func codeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    renderTemplate(w, "code", &struct{Ref string}{Ref: ps.ByName("ref")})
}

func generateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	addr := ps.ByName("address")
	ref := ps.ByName("ref");
	fmt.Println("Addr: %s, Ref: %s", addr, ref)
	code := tracking_code.GenerateCode(12)
	renderTemplate(w, "generate", &struct{Code string}{Code: code})
}

func myInvitationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderTemplate(w, "my-invitation", &struct{Code string}{Code: "9527"})
}


var (
	templateMap = template.FuncMap{
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}
	templates   = template.New("").Funcs(templateMap)
	templateBox *rice.Box
)

func newTemplate(path string, _ os.FileInfo, _ error) error {
	if path == "" {
		return nil
	}
	templateString, err := templateBox.String(path)
	if err != nil {
		log.Panicf("Unable to extract: path=%s, err=%s", path, err)
	}
	if _, err = templates.New(filepath.Join("tpl", path)).Parse(templateString); err != nil {
		log.Panicf("Unable to parse: path=%s, err=%s", path, err)
	}
	return nil
}

// Render a template given a model
func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, "tpl/"+tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func getRouter() *httprouter.Router {
	// Load and parse templates (from binary or disk)
	templateBox = rice.MustFindBox("tpl")
	templateBox.Walk("", newTemplate)

	// mux handler
	router := httprouter.New()

	// Index routee
	router.GET("/code/", codeHandler)

	// Example route that takes one rest style option
	router.POST("/code/generate/", generateHandler)

	// Example route that encounters an error
	router.POST("/code/my-invitation/", myInvitationHandler)

	// Serve static assets via the "static" directory
	fs := rice.MustFindBox("s").HTTPBox()
	router.ServeFiles("/s/*filepath", fs)
	return router
}