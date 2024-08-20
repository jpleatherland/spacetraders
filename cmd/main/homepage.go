package main

import (
	"net/http"
	"text/template"

	"github.com/jpleatherland/spacetraders/internal/db"
)

type ServerStatus struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	ResetDate   string `json:"resetDate"`
	Description string `json:"description"`
	Stats       struct {
		Agents    int `json:"agents"`
		Ships     int `json:"ships"`
		Systems   int `json:"systems"`
		Waypoints int `json:"waypoints"`
	} `json:"stats"`
	Leaderboards struct {
		MostCredits []struct {
			AgentSymbol string `json:"agentSymbol"`
			Credits     int64  `json:"credits"`
		} `json:"mostCredits"`
		MostSubmittedCharts []struct {
			AgentSymbol string `json:"agentSymbol"`
			ChartCount  int    `json:"chartCount"`
		} `json:"mostSubmittedCharts"`
	} `json:"leaderboards"`
	ServerResets struct {
		Next      string `json:"next"`
		Frequency string `json:"frequency"`
	} `json:"serverResets"`
	Announcements []struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"announcements"`
	Links []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"links"`
}

func (resources *Resources) homePage(rw http.ResponseWriter, req *http.Request, session db.Session) {
	tmpl := template.Must(template.ParseFiles("views/homepage.html"))

	resp := http.Get("")

	err := tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, "error rendering template", http.StatusInternalServerError)
		return
	}
}
