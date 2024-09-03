package routes

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
)

func SetSession(rw http.ResponseWriter, req *http.Request) {
	symbol := req.PathValue("agentSymbol")

	resources, ok := middleware.GetResources(req.Context())
	if !ok {
		log.Println("unable to get resources in set session")
		response.RespondWithHTMLError(rw, "failed to set agent session", http.StatusInternalServerError)

	}

	session, ok := middleware.GetSession(req.Context())

	agents, err := resources.DB.GetAgentsByUserId(req.Context(), session.UserID)
	if err != nil {
		log.Println("failed to get agents by user id: ", err.Error())
		response.RespondWithHTMLError(rw, "server error", http.StatusInternalServerError)
		return
	}

	sessionAgentParams := db.SetAgentForSessionParams{ID: session.ID, AgentID: uuid.NullUUID{}}
	for _, agent := range agents {
		if agent.Name == symbol {
			sessionAgentParams.AgentID.UUID = agent.ID
			sessionAgentParams.AgentID.Valid = true
			break
		}
	}

	if !sessionAgentParams.AgentID.Valid {
		response.RespondWithHTMLError(rw, "unable to find agent for session", http.StatusInternalServerError)
		return
	}

	err = resources.DB.SetAgentForSession(req.Context(), sessionAgentParams)
	if err != nil {
		response.RespondWithHTMLError(rw, "unable to set agent for session", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("HX-Redirect", "http://localhost:8080/game")
	rw.WriteHeader(http.StatusOK)
}
