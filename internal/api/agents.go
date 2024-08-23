package api

import (
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

// List Agents
// (GET /agents)
func (s Server) GetAgents(w http.ResponseWriter, r *http.Request, params spec.GetAgentsParams) {
	routes.GetAgentTest()
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
