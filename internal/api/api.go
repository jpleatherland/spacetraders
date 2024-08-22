package api

import (
	"net/http"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) GetStatus(rw http.ResponseWriter, req *http.Request) {

}
func (s Server) GetAgents(rw http.ResponseWriter, req *http.Request, params GetAgentsParams) {

}
