package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/db"
)

func (s Server) GetAgentHandler(database *db.Queries, agentId uuid.UUID) (Agent, error) {
	agent := Agent{}
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "", nil)
	req.Header.Set("Authorization", "Bearer "+agentId.String())
	s.GetMyAgent(rw, req)
	if rw.Result().StatusCode != http.StatusOK {
		return agent, errors.New(strconv.Itoa(rw.Result().StatusCode))
	}
	body, err := io.ReadAll(rw.Result().Body)
	if err != nil {
		return agent, err
	}

	err = json.Unmarshal(body, &agent)
	if err != nil {
		return agent, err
	}

	return agent, nil
}

// List Agents
// (GET /agents)
func (s Server) GetAgents(w http.ResponseWriter, r *http.Request, params GetAgentsParams) {}

// Get Public Agent
// (GET /agents/{agentSymbol})
func (s Server) GetAgent(w http.ResponseWriter, r *http.Request, agentSymbol string) {}

// Get Agent
// (GET /my/agent)
func (s Server) GetMyAgent(w http.ResponseWriter, r *http.Request) {

}
