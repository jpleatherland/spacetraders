package api

import (
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/routes"
)

func (s Server) GetStatus(w http.ResponseWriter, r *http.Request) {
	res, err := routes.GetStatusHandler()
	if err != nil {
		response.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, res)
}
