package main

import (
	"fmt"
	"html/template"
	"net/http"

)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR: %s", err)
		}
	}()
	http.HandleFunc("/code/", codeHandler)
    http.HandleFunc("/code/generate/", generateHandler)
    http.HandleFunc("/code/my-invitation/", myInvitationHandler)
    err := http.ListenAndServe(":8888", nil)
    if err != nil {
        println("error： %s", err)
    }
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ref := r.FormValue("ref")
    renderTemplate(w, "code", struct{Ref string}{Ref: ref})
}

func generateHandler(w http.ResponseWriter, r *http.Request) {

}

func myInvitationHandler(w http.ResponseWriter, r *http.Request) {

}

var templates = template.Must(template.ParseFiles("tpl/code.html", "tpl/generate.html","tpl/my-invitation.html"))


func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, "tpl/"+tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
