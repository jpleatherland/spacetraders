package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

func CreateAgent(req *http.Request) (spec.RegisterResponse, int, error) {
	url, ok := middleware.GetUrlContext(req.Context())
	if !ok {
		return spec.RegisterResponse{}, http.StatusInternalServerError, errors.New("unable to get register URL")
	}

	registerResult := spec.RegisterResponse{}

	body, _ := io.ReadAll(req.Body)

	resp, err := http.Post(baseUrl+url, "application/json", bytes.NewReader(body))

	if err != nil {
		return registerResult, resp.StatusCode, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return registerResult, resp.StatusCode, err
	}

	if resp.StatusCode > 299 {
		errMsg := fmt.Sprintf("Unexpected status code received: %v", string(bytes[:]))
		return registerResult, resp.StatusCode, errors.New(errMsg)
	}

	err = json.Unmarshal([]byte(bytes), &registerResult)
	return registerResult, resp.StatusCode, err
}

func GetAgents(resources *middleware.Resources, session db.Session) ([]spec.Agent, error) {
	ctx := context.Background()
	agents := []spec.Agent{}
	agentIds, err := resources.DB.GetAgentsByUserId(ctx, session.UserID)
	if err != nil {
		return agents, err
	}

	if len(agentIds) == 0 {
		return agents, errors.New("no agents found")
	}

	wg := sync.WaitGroup{}
	results := make(chan spec.Agent, len(agentIds))
	wg.Add(len(agentIds))

	for i := range agentIds {
		go func() {
			defer wg.Done()
			cacheAgent, ok := resources.Cache.Get(agentIds[i].ID.String())
			if ok {
				log.Printf("found %v in cache\n", agentIds[i].ID.String())
				t, ok := cacheAgent.(spec.Agent)
				if ok {
					results <- t
					return
				}
			}
			agent, err := GetAgentHandler(resources, agentIds[i].ID)
			if err != nil {
				log.Printf("Get agent failed for %v: %v", agentIds[i].ID.String(), err.Error())
				return
			}
			writeAgentToCache(resources, agentIds[i].ID.String(), agent)
			results <- agent
		}()
	}

	wg.Wait()
	close(results)

	chIdx := 0
	for agent := range results {
		agents = append(agents, agent)
		chIdx++
	}
	return agents, nil
}

func GetAgentHandler(resources *middleware.Resources, agentId uuid.UUID) (spec.Agent, error) {
	agent := spec.AgentResponse{}
	myAgentUrl := "/my/agent"

	req, err := http.NewRequest("GET", baseUrl+myAgentUrl, nil)
	if err != nil {
		return agent.Data, err
	}
	ctx := context.Background()

	agentToken, err := resources.DB.GetAgentTokenById(ctx, agentId)
	if err != nil {
		return agent.Data, err
	}

	req.Header.Set("Authorization", "Bearer "+agentToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return agent.Data, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return agent.Data, err
	}

	if res.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("%v: %v", res.StatusCode, body)
		return agent.Data, errors.New(errMsg)
	}

	err = json.Unmarshal(body, &agent)
	if err != nil {
		log.Println(err.Error())
	}

	return agent.Data, nil
}

func writeAgentToCache(resources *middleware.Resources, agentId string, agent spec.Agent) error {
	log.Printf("adding %v to cache\n", agentId)
	resources.Cache.Add(
		agentId,
		agent,
		time.Duration(5*time.Minute),
	)
	return nil
}

func writeAgentToDB(resources *middleware.Resources, registerResp spec.RegisterResponse, userId uuid.UUID) error {
	dbAgentParams := db.CreateAgentParams{
		ID:     uuid.New(),
		Name:   registerResp.Data.Agent.Symbol,
		Token:  registerResp.Data.Token,
		UserID: userId,
	}
	err := resources.DB.CreateAgent(context.Background(), dbAgentParams)
	return err
}

func RegisterAgent(rw http.ResponseWriter, req *http.Request) {
	resources, ok := middleware.GetResources(req.Context())
	if !ok {
		response.RespondWithHTMLError(rw, "Request failure", http.StatusInternalServerError)
		return
	}

	session, ok := middleware.GetSession(req.Context())
	if !ok {
		response.RespondWithHTMLError(rw, "Request failure", http.StatusInternalServerError)
		return
	}

	registeredAgent, respCode, err := CreateAgent(req)
	if err != nil && respCode == http.StatusCreated {
		htmlResp := fmt.Sprintf("Agent registered but the response wasn't handled: %v", err.Error())
		response.RespondWithHTMLError(rw, htmlResp, http.StatusInternalServerError)
		return
	}
	if err != nil && respCode != http.StatusCreated {
		htmlResp := fmt.Sprintf("Agent not registered: %v", err.Error())
		response.RespondWithHTMLError(rw, htmlResp, http.StatusInternalServerError)
		return
	}
	if respCode != http.StatusCreated {
		htmlResp := fmt.Sprintf("Agent not registered: %v", err.Error())
		response.RespondWithHTMLError(rw, htmlResp, http.StatusInternalServerError)
	}

	err = writeAgentToDB(resources, registeredAgent, session.UserID)
	if err != nil {
		response.RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
	}

	response.RespondWithHTML(rw, "Agent registered. Select from the list to continue with this agent", respCode)
}
