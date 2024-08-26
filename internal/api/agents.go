package api

import (
	"log"
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

// List Agents
// (GET /agents)
func (s Server) GetAgents(w http.ResponseWriter, r *http.Request, params spec.GetAgentsParams) {
	routes.GetAgentTest()
	log.Println("in get agents")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hi"))
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
func Register(w http.ResponseWriter, r *http.Request) {
	routes.RegisterAgent(w, r)
}
