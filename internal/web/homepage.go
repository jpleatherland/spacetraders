package web

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

type HomePageData struct {
	Server spec.ServerStatus
	Agents []db.GetAgentsByUserIdRow
}

func HomePage(rw http.ResponseWriter, req *http.Request) {
	log.Println("in homepage")
	resources, ok := middleware.GetResources(req.Context())
	if !ok {
		response.RespondWithError(rw, "unable to load resources", http.StatusInternalServerError)
	}

	_, ok = middleware.GetSession(req.Context())
	if !ok {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return
	}

	tmpl, err := template.Must(template.ParseGlob(filepath.Join("views", "templates", "homepage.html"))).ParseGlob(filepath.Join("views", "templates", "partials", "*.html"))
	if err != nil {
		response.RespondWithError(rw, "Unable to load templates", http.StatusInternalServerError)
	}

	pageData := HomePageData{
		Server: spec.ServerStatus{},
		Agents: []db.GetAgentsByUserIdRow{},
	}

	statusResp, exists := resources.Cache.Get("serverStatus")

	if !exists {
		statusResp, err = routes.GetStatusHandler()
		if err != nil {
			pageData.Server.RequestStatus = err.Error()
			response.RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		resources.Cache.Add("serverStatus", statusResp, 0)
	}

	result, ok := statusResp.(spec.ServerStatus)
	if !ok {
		log.Print("cached status is not of type ServerStatus")
	} else {
		pageData.Server = result
	}

	//pageData.Agents = routes.GetAgents(resources, session.UserID)

	err = tmpl.Execute(rw, pageData)
	if err != nil {
		log.Printf("error rendering template: %v", err.Error())
	}
}
