package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

func GetWaypoint(rw http.ResponseWriter, req *http.Request, systemSymbol, waypointSymbol string) {
	log.Println("in get waypoint")
	resources, ok := middleware.GetResources(req.Context())
	if !ok {
		response.RespondWithError(rw, "couldn't get resources", http.StatusInternalServerError)
		return
	}

	cachedWaypoint, ok := resources.Cache.Get(waypointSymbol)
	if ok {
		response.RespondWithPartialTemplate(rw, "systemPartials", "waypoint.html", cachedWaypoint)
	}

	log.Println("not found in cache, getting waypoint info")
	url, err := url.JoinPath(baseUrl, "systems", systemSymbol, "waypoints", waypointSymbol)
	if err != nil {
		log.Println("unable to create waypoint url: ", err.Error())
		response.RespondWithHTMLError(rw, "unable to get waypoint information", http.StatusOK)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTMLError(rw, "unable to get waypoint information", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to unmarshal waypoint information: ", err.Error())
		response.RespondWithHTMLError(rw, "unable to get waypoint information", http.StatusInternalServerError)
		return
	}

	waypointResponse := spec.WaypointResponse{}
	err = json.Unmarshal(body, &waypointResponse)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	resources.Cache.Add(waypointSymbol, waypointResponse, 0)
	response.RespondWithPartialTemplate(rw, "systemPartials", "waypoint.html", waypointResponse)

}
