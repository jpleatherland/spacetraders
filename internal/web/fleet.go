package web

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/response"
)

func Fleet(rw http.ResponseWriter, req *http.Request) {

	tmpl := template.Must(template.ParseGlob(filepath.Join("views", "templates", "fleet.html")))
	err := tmpl.Execute(rw, nil)
	if err != nil {
		response.RespondWithHTML(rw, "Unable to load systems", http.StatusInternalServerError)
	}
}
