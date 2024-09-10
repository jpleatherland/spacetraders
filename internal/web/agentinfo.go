package web

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/spec"
)


func AgentInfo(rw http.ResponseWriter, req *http.Request) {
	resources, ok := middleware.GetResources(req.Context())
	if !ok {
		response.RespondWithError(rw, "unable to load resources", http.StatusInternalServerError)
	}
	session, ok := middleware.GetSession(req.Context())
	if !ok {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return
	}
	pageData := spec.Agent{}
	agent, err := routes.GetAgentHandler(resources, session.AgentID.UUID)
	if err != nil {
		response.RespondWithHTML(rw, "unable to load agents", http.StatusInternalServerError)
	}
	pageData = agent
	tmpl := template.Must(template.ParseGlob(filepath.Join("views", "templates", "agentinfo.html")))
	err = tmpl.Execute(rw, pageData)
	if err != nil {
		response.RespondWithHTML(rw, "Unable to load systems", http.StatusInternalServerError)
	}
}
