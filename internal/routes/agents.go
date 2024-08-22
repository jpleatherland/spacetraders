package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/api"
	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/resources"
)

func createAgent(username, faction string) (string, error) {
	createAgentUrl := "https://api.spacetraders.io/v2/register"

	requestBody := struct {
		Symbol  string `json:"symbol"`
		Faction string `json:"faction"`
	}{
		Symbol:  username,
		Faction: faction,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(jsonBody)

	resp, err := http.Post(createAgentUrl, "application/json", reader)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &result)
	if err != nil {
		return "", err
	}

	return result["data"].(map[string]interface{})["token"].(string), nil
}

func GetAgents(database *db.Queries, server api.Server, userId uuid.UUID) ([]db.GetAgentsByUserIdRow, error) {
	ctx := context.Background()
	agents := []db.GetAgentsByUserIdRow{}
	agents, err := database.GetAgentsByUserId(ctx, userId)
	if err != nil {
		return agents, err
	}

	if len(agents) == 0 {
		return agents, nil
	}

	wg := sync.WaitGroup{}
	results := make(chan api.Agent)
	wg.Add(len(agents))

	for i := range agents {
		go func() {
			defer wg.Done()
			result, err := server.GetAgentHandler(database, agents[i].ID)
			if err != nil {
				log.Printf("failed to get agent")
			}
			results <- result
		}()
	}
	// this agents need to be augmented
	// by ranging over results 
	// which is a call to st api getting agent details
	return agents, nil
}
