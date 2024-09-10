package web

import (
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/response"
)

func Contracts(rw http.ResponseWriter, _ *http.Request) {
	response.RespondWithTemplate(rw, "contracts.html", nil)
}
