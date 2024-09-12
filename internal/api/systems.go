package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

func (s Server) GetWaypoint(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
	log.Println("in get waypoint")
	resources, ok := middleware.GetResources(r.Context())
	if !ok {
		response.RespondWithError(w, "couldn't get resources", http.StatusInternalServerError)
		return
	}

	cachedWaypoint, ok := resources.Cache.Get(waypointSymbol)
	if ok {
		response.RespondWithPartialTemplate(w, "systemPartials", "waypoint.html", cachedWaypoint)
		return
	}

	log.Println("not found in cache, getting waypoint info")
	url, err := url.JoinPath(baseUrl, "systems", systemSymbol, "waypoints", waypointSymbol)
	if err != nil {
		log.Println("unable to create waypoint url: ", err.Error())
		response.RespondWithHTMLError(w, "unable to get waypoint information", http.StatusOK)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTMLError(w, "unable to get waypoint information", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to unmarshal waypoint information: ", err.Error())
		response.RespondWithHTMLError(w, "unable to get waypoint information", http.StatusInternalServerError)
		return
	}

	waypointResponse := spec.WaypointResponse{}
	err = json.Unmarshal(body, &waypointResponse)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTMLError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resources.Cache.Add(waypointSymbol, waypointResponse, 0)
	response.RespondWithPartialTemplate(w, "systemPartials", "waypoint.html", waypointResponse)
}

func (s Server) GetSystem(w http.ResponseWriter, r *http.Request, systemSymbol string) {
	log.Println("in get system")
}

func (s Server) GetSystemWaypoints(w http.ResponseWriter, r *http.Request, systemSymbol string, params spec.GetSystemWaypointsParams) {
	log.Println("in get systemwaypoints")
}

func (s Server) GetSystems(w http.ResponseWriter, r *http.Request, params spec.GetSystemsParams) {
	resources, ok := middleware.GetResources(r.Context())
	if !ok {
		response.RespondWithError(w, "couldn't get resources", http.StatusInternalServerError)
		return
	}

	page := "1"
	if params.Page != nil {
		page = strconv.Itoa(*params.Page)
	}

	limit := "10"
	if params.Limit != nil {
		limit = strconv.Itoa(*params.Limit)
	}

	cachedWaypoint, ok := resources.Cache.Get("systemList"+page+limit)
	if ok {
		response.RespondWithPartialTemplate(w, "systemPartials", "waypoint.html", cachedWaypoint)
		return
	}

	log.Println("not found in cache, getting systems")

	path := "/systems"
	queryParams := url.Values{}
	queryParams.Add("page", page)
	queryParams.Add("limit", limit)

	u, _ := url.ParseRequestURI(baseUrl)
	u.Path = path
	u.RawQuery = queryParams.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("unable to get systems: ", err.Error())
		response.RespondWithHTML(w, "unable to get systems", http.StatusOK)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to unmarshal system list: ", err.Error())
		response.RespondWithHTMLError(w, "unable to get system list", http.StatusOK)
		return
	}

	systemListResponse := spec.SystemList{}
	err = json.Unmarshal(body, &systemListResponse)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTMLError(w, err.Error(), http.StatusOK)
		return
	}

	resources.Cache.Add("systemList"+page+limit, systemListResponse, 120)
	response.RespondWithPartialTemplate(w, "systemPartial", "systemList.html", systemListResponse)
}
