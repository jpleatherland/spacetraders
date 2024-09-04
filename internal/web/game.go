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

type GamePageData struct {
	Agent spec.Agent
}

func Game(rw http.ResponseWriter, req *http.Request) {
	session, ok := middleware.GetSession(req.Context())
	if !ok {
		log.Println("can't get session")
	}

	resources, ok := middleware.GetResources(req.Context())
	if !ok {
		log.Println("can't get agent")
	}

	agent, err := routes.GetAgentHandler(resources, session.AgentID.UUID)

	tmpl, err := template.Must(template.ParseGlob(filepath.Join("views", "templates", "game.html"))).ParseGlob(filepath.Join("views", "templates", "partials", "*.html"))
	if err != nil {
		response.RespondWithError(rw, "Unable to load templates", http.StatusInternalServerError)
	}

	pageData := GamePageData{Agent: agent}

	err = tmpl.Execute(rw, pageData)
	if err != nil {
		log.Printf("error rendering template: %v", err.Error())
	}
}
