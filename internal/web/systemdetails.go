package web

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/response"
)

func SystemDetails(rw http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseGlob(filepath.Join("views", "templates", "systemdetails.html")))
	err := tmpl.Execute(rw, nil)
	if err != nil {
		log.Println("error ", err.Error())
		response.RespondWithHTML(rw, "Unable to load systems", http.StatusInternalServerError)
	}
}
