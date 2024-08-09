package main

import (
	"net/http"
	"text/template"
)

func (apiConf *apiConfig) loginPage(rw http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/login.html"))

	err := tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, "error rendering template", http.StatusInternalServerError)
		return
	}
}
