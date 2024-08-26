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
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

func CreateAgent(username, faction string) (string, error) {
	createAgentUrl := "/register"

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

func GetAgents(resources *middleware.Resources, userId uuid.UUID) ([]spec.Agent, error) {
	ctx := context.Background()
	agents := []spec.Agent{}
	agentIds, err := resources.DB.GetAgentsByUserId(ctx, userId)
	if err != nil {
		return agents, err
	}

	if len(agents) == 0 {
		return agents, errors.New("no agents found")
	}

	wg := sync.WaitGroup{}
	results := make(chan spec.Agent)
	wg.Add(len(agents))

	for i := range agentIds {
		go func() {
			defer wg.Done()
			agent, err := GetAgentHandler(agentIds[i].ID)
			if err != nil {
				log.Printf("Get agent failed for %v: %v", agentIds[i].ID.String(), err.Error())
			}
			results <- agent
		}()
	}
	for agent := range results {
		agents = append(agents, agent)
		writeAgentToCache(resources, agent)
	}
	return agents, nil
}

func GetAgentHandler(agentId uuid.UUID) (spec.Agent, error) {
	agent := spec.Agent{}
	myAgentUrl := "/my/agent"

	req, err := http.NewRequest("GET", baseUrl+myAgentUrl, nil)
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

func writeAgentToCache(resources *middleware.Resources, agent spec.Agent) {
	resources.Cache.Add(
		*agent.AccountId,
		agent,
		time.Duration(5*time.Minute),
	)
}

func RegisterAgent(rw http.ResponseWriter, req *http.Request) {
	response.RespondWithHTML(rw, "<p>Agent register hit. Not registered.</p>", http.StatusOK)
}
