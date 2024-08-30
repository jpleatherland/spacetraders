package routes

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

func SetSession(rw http.ResponseWriter, req *http.Request) {
	symbol := strings.TrimPrefix(req.URL.Path, "/products/")
	log.Println("in set session", symbol)
	time.Sleep(5 * time.Second)
	resources, ok := middleware.GetResources(req.Context())
	if !ok {
		log.Println("unable to get resources in set session")
		response.RespondWithHTMLError(rw, "failed to set agent session", http.StatusInternalServerError)

	}
	statusResp, exists := resources.Cache.Get("serverStatus")

	if !exists {
		statusResp, _ = GetStatusHandler(resources)
	}

	status, ok := statusResp.(spec.ServerStatus)
	if !ok {
		log.Println("server status not of type ServerStatus, unable to set reset date on agent")
	}
	
	resetDate := status.ServerResets.Next	

	
}
