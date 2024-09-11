package web

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/routes"
)

func SystemDetails(rw http.ResponseWriter, req *http.Request) {
	system := req.PathValue(system)
	waypoint := req.PathValue(waypoint)
	routes.GetWaypoint(rw, req, system, waypoint)
	tmpl := template.Must(template.ParseGlob(filepath.Join("views", "templates", "systemdetails.html")))
	err := tmpl.Execute(rw, nil)
	if err != nil {
		response.RespondWithHTML(rw, "Unable to load systems", http.StatusInternalServerError)
	}
}
