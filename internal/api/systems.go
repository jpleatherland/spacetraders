package api

import (
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/spec"
)


func (s Server) GetWaypoint(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
	routes.GetWaypoint(w, r, systemSymbol, waypointSymbol)
}

func (s Server) GetSystem(w http.ResponseWriter, r *http.Request, systemSymbol string)            {}

func (s Server) GetSystemWaypoints(w http.ResponseWriter, r *http.Request, systemSymbol string, params spec.GetSystemWaypointsParams) {
}

func (s Server) GetSystems(w http.ResponseWriter, r *http.Request, params spec.GetSystemsParams)  {}
