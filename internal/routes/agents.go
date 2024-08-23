package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/spec"
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

func GetAgents(database *db.Queries, userId uuid.UUID) ([]db.GetAgentsByUserIdRow, error) {
	ctx := context.Background()
	agents := []db.GetAgentsByUserIdRow{}
	agents, err := database.GetAgentsByUserId(ctx, userId)
	if err != nil {
		return agents, err
	}

	if len(agents) == 0 {
		return agents, nil
	}

	// wg := sync.WaitGroup{}
	// results := make(chan api.Agent)
	// wg.Add(len(agents))

	// for i := range agents {
	// 	go func() {
	// 		defer wg.Done()
	// 		// result, err := server.GetAgentHandler(database, agents[i].ID)
	// 		// if err != nil {
	// 			// log.Printf("failed to get agent")
	// 		// }
	// 		results <- result
	// 	}()
	// }
	// this agents need to be augmented
	// by ranging over results
	// which is a call to st api getting agent details
	return agents, nil
}

func GetAgentHandler(database *db.Queries, agentId uuid.UUID) (spec.Agent, error) {
	agent := spec.Agent{}

	req, err := http.NewRequest("GET", baseUrl+"/my/agent", nil)
	if err != nil {
		return agent, err
	}
	req.Header.Set("Authorization", "Bearer "+agentId.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return agent, err
	}

	if res.StatusCode != http.StatusOK {
		return agent, errors.New(strconv.Itoa(res.StatusCode))
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return agent, err
	}

	err = json.Unmarshal(body, &agent)
	if err != nil {
		return agent, err
	}

	return agent, nil
}

func GetAgentTest() {
	log.Println("in routes")
}
