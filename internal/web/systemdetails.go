package web

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/response"
)

type SystemDetailsPageData struct {
	System      string
	Waypoint    string
}

func SystemDetails(rw http.ResponseWriter, req *http.Request) {
	system := req.PathValue("system")
	waypoint := req.PathValue("waypoint")
	pageData := SystemDetailsPageData{
		System:   system,
		Waypoint: waypoint,
	}

	tmpl := template.Must(template.ParseGlob(filepath.Join("views", "templates", "systemdetails.html")))
	err := tmpl.Execute(rw, pageData)
	if err != nil {
		log.Println("error ", err.Error())
		response.RespondWithHTML(rw, "Unable to load systems", http.StatusInternalServerError)
	}
}
