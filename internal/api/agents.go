package api

import (
	"log"
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

// List Agents
// (GET /agents)
func (s Server) GetAgents(w http.ResponseWriter, r *http.Request, params spec.GetAgentsParams) {
}

// Get Public Agent
// (GET /agents/{agentSymbol})
func (s Server) GetAgent(w http.ResponseWriter, r *http.Request, agentSymbol string) {}

// Get Agent
// (GET /my/agent)
func (s Server) GetMyAgent(w http.ResponseWriter, r *http.Request) {

}

// Register New Agent
// (POST /register)
func (s Server) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("In register")
	reqWithContext := middleware.UrlContext("/register", r)
	log.Println(r.Context().Value(middleware.UrlKey))
	log.Println(reqWithContext.Context().Value(middleware.UrlKey))
	routes.RegisterAgent(w, reqWithContext)
}
