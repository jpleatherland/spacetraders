package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
)

func SetSession(rw http.ResponseWriter, req *http.Request) {
	symbol := req.PathValue("symbol")
	log.Println("in set session", symbol)
	time.Sleep(5 * time.Second)

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
	}
	for _, agent := range agents {
		if agent.Name == symbol {
			session.AgentID = uuid.NullUUID{
				UUID: agent.ID,
				Valid: true,
			}
			break
		}
	}

}
