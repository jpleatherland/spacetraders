package web

import (
	"log"
	"net/http"
	"github.com/jpleatherland/spacetraders/internal/middleware"
)

func Game(rw http.ResponseWriter, req *http.Request) {
	session, ok := middleware.GetSession(req.Context())
	if !ok {
		log.Println("can't get session")
	}
	log.Println(session.AgentID.UUID)
}
