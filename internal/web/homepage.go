package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"text/template"
	"path/filepath"

	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/routes"
)

type ServerStatus struct {
	ServerResets struct {
		Next      string `json:"next"`
		Frequency string `json:"frequency"`
	} `json:"serverResets"`
	RequestStatus string `json:"RequestStatus"`
	Status        string `json:"status"`
	Version       string `json:"version"`
	ResetDate     string `json:"resetDate"`
	Description   string `json:"description"`
	Leaderboards  struct {
		MostCredits []struct {
			AgentSymbol string `json:"agentSymbol"`
			Credits     int64  `json:"credits"`
		} `json:"mostCredits"`
		MostSubmittedCharts []struct {
			AgentSymbol string `json:"agentSymbol"`
			ChartCount  int    `json:"chartCount"`
		} `json:"mostSubmittedCharts"`
	} `json:"leaderboards"`
	Announcements []struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"announcements"`
	Links []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"links"`
	Stats struct {
		Agents    int `json:"agents"`
		Ships     int `json:"ships"`
		Systems   int `json:"systems"`
		Waypoints int `json:"waypoints"`
	} `json:"stats"`
}

type HomePageData struct {
	Server ServerStatus
	Agents []db.GetAgentsByUserIdRow
}

func HomePage(rw http.ResponseWriter, _ *http.Request, session db.Session, resources *routes.Resources) {
	tmpl, err := template.Must(template.ParseGlob(filepath.Join("views", "templates", "homepage.html"))).ParseGlob(filepath.Join("views", "templates", "partials", "*.html"))
	//tmpl := template.Must(template.ParseFiles(
	//	"views/templates/homepage.html",
	//	"views/templates/partials/agents.html",
	//	"views/templates/partials/leaderboards.html",
	//))
	//tmpl := template.Must(template.ParseFiles("views/homepage.html").template.ParseFiles())
	serverStatus := ServerStatus{}
	log.Printf("session id: %v", session.ID)

	body, exists := resources.Cache.Get("serverStatus")
	if !exists {
		resp, err := http.Get("https://api.spacetraders.io/v2/")
		if err != nil {
			serverStatus.RequestStatus = err.Error()
		}

		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			response.RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		resources.Cache.Add("serverStatus", body, 0)
	}

	serverStatus.RequestStatus = "OK"
	err = json.Unmarshal(body, &serverStatus)
	if err != nil {
		log.Printf("unable to unmarshal response: %v", err.Error())
	}

	pageData := HomePageData{
		Server: serverStatus,
		Agents: []db.GetAgentsByUserIdRow{},
	}	

	err = tmpl.Execute(rw, pageData)
	if err != nil {
		log.Printf("error rendering template: %v", err.Error())
	}
}
