package api

import (
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/routes"
)

func (s Server) GetStatus(w http.ResponseWriter, r *http.Request) {
	resources, ok := middleware.GetResources(r.Context())
	if !ok {
		response.RespondWithError(w, "unable to fetch server resources", http.StatusInternalServerError)
	}
	res, err := routes.GetStatusHandler(resources)
	if err != nil {
		response.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, res)
}
