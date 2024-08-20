package main

import (
	"net/http"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/db"
)

func (resources *Resources) homePage(rw http.ResponseWriter, req *http.Request, session db.Session) {
	tmpl := template.Must(template.ParseFiles("views/homepage.html"))

	err := tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, "error rendering template", http.StatusInternalServerError)
		return
	}
}
