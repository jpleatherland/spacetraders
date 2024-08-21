package web

import (
	"net/http"
	"text/template"
)

func LoginPage(rw http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/login.html"))

	err := tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, "error rendering template", http.StatusInternalServerError)
		return
	}
}
