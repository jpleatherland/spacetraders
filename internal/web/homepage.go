package web

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/api"
	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/routes"
)

type HomePageData struct {
	Server api.ServerStatus
	Agents []db.GetAgentsByUserIdRow
}

func HomePage(rw http.ResponseWriter, req *http.Request, session db.Session, resources *routes.Resources) {
	tmpl, err := template.Must(template.ParseGlob(filepath.Join("views", "templates", "homepage.html"))).ParseGlob(filepath.Join("views", "templates", "partials", "*.html"))
	if err != nil {
		response.RespondWithError(rw, "Unable to load templates", http.StatusInternalServerError)
	}

	pageData := HomePageData{
		Server: api.ServerStatus{},
		Agents: []db.GetAgentsByUserIdRow{},
	}

	statusResp, exists := resources.Cache.Get("serverStatus")

	if !exists {
		statusResp, err = resources.Server.GetStatusHandler(rw, req)
		if err != nil {
			pageData.Server.RequestStatus = err.Error()
			response.RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		resources.Cache.Add("serverStatus", statusResp, 0)
	}

	result, ok := statusResp.(api.ServerStatus)
	if !ok {
		log.Print("cached status is not of type ServerStatus")
	} else {
		pageData.Server = result
	}

	err = tmpl.Execute(rw, pageData)
	if err != nil {
		log.Printf("error rendering template: %v", err.Error())
	}
}
