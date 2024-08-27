package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
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

	resp, err := http.Post(baseUrl+url, "application/json", req.Body)

	if err != nil {
		return registerResult, resp.StatusCode, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return registerResult, resp.StatusCode, err
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
		writeAgentToCache(resources, agent, session)
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

func writeAgentToCache(resources *middleware.Resources, agent spec.Agent, session db.Session) error {
	if !session.AgentID.Valid {
		return errors.New("agentId is not valid, unable to cache")
	}
	resources.Cache.Add(
		session.AgentID.UUID.String(),
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

	err = writeAgentToDB(resources, registeredAgent, session.UserID)
	if err != nil {
		response.RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
	}

	response.RespondWithHTML(rw, "Agent registered. Select from the list to continue with this agent", respCode)
}
