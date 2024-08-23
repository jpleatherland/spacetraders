package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/spec"
)

func GetStatusHandler() (spec.ServerStatus, error) {
	statusResp := spec.ServerStatus{}
	resp, err := http.Get(baseUrl + "/")
	if err != nil {
		return statusResp, err
	}
	if resp.StatusCode != http.StatusOK {
		statusResp.RequestStatus = "request error"
		return statusResp, errors.New("error fetching server status")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		statusResp.RequestStatus = "request error"
		return statusResp, err
	}

	err = json.Unmarshal(body, &statusResp)
	if err != nil {
		statusResp.RequestStatus = "request error"
		return statusResp, err
	}
	statusResp.RequestStatus = "ok"
	return statusResp, nil
}
