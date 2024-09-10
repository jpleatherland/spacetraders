package web

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

type homePageData struct {
	Server        spec.ServerStatus
	AgentsSection struct {
		Agents   []spec.Agent
		Factions [19]spec.FactionSymbol
	}
}

func HomePage(rw http.ResponseWriter, req *http.Request) {
	resources, ok := middleware.GetResources(req.Context())
	if !ok {
		response.RespondWithError(rw, "unable to load resources", http.StatusInternalServerError)
	}

	session, ok := middleware.GetSession(req.Context())
	if !ok {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return
	}

	tmpl, err := template.Must(template.ParseGlob(filepath.Join("views", "templates", "homepage.html"))).ParseGlob(filepath.Join("views", "templates", "partials", "*.html"))
	if err != nil {
		log.Println(err.Error())
		response.RespondWithError(rw, "Unable to load templates", http.StatusInternalServerError)
	}

	pageData := homePageData{}

	statusResp, exists := resources.Cache.Get("serverStatus")

	if !exists {
		statusResp, err = routes.GetStatusHandler(resources)
		if err != nil {
			log.Println(err.Error())
			pageData.Server.RequestStatus = "Unable to get server status"
		}
	}

	result, ok := statusResp.(spec.ServerStatus)
	if !ok {
		log.Print("cached status is not of type ServerStatus")
	} else {
		pageData.Server = result
	}

	pageData.AgentsSection.Agents, err = routes.GetAgents(resources, session)
	if err != nil {
		log.Println(err.Error())
	}
	pageData.AgentsSection.Factions = spec.Factions

	err = tmpl.Execute(rw, pageData)
	if err != nil {
		log.Printf("error rendering template: %v", err.Error())
	}
}
