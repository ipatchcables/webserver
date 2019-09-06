package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/process", process)
	log.Fatal(http.ListenAndServe("10.9.1.5:80", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func process(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user := r.FormValue("nameBox")
	pass := r.FormValue("pwdBox")

	d := struct {
		Username    string
		Supersecret string
	}{
		Username:    user,
		Supersecret: pass,
	}
	tpl.ExecuteTemplate(w, "process.html", d)

}
