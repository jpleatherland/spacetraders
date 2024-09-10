package web

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/response"
)

func Contracts(rw http.ResponseWriter, _ *http.Request){
	tmpl := template.Must(template.ParseGlob(filepath.Join("views", "templates", "contracts.html")))
	err := tmpl.Execute(rw, nil)
	if err != nil {
		response.RespondWithHTML(rw, "Unable to load contracts", http.StatusInternalServerError)
	}
}

