package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
)

func (s Server) GetStatus(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(BASEURL)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s Server) GetStatusHandler(w http.ResponseWriter, r *http.Request) (ServerStatus, error) {
	statusResp := ServerStatus{}
	rw := httptest.NewRecorder()
	s.GetStatus(rw, r)
	if rw.Result().StatusCode != http.StatusOK {
		statusResp.RequestStatus = "request error"
		return statusResp, errors.New("Error fetching server status")
	}
	body, err := io.ReadAll(rw.Result().Body)
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
