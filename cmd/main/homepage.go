package main

import (
	"net/http"
	"text/template"
)

func (apiConf *apiConfig) homePage(rw http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/homepage.html"))

	err := tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, "error rendering template", http.StatusInternalServerError)
		return
	}
}
