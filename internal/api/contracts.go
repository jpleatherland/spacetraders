package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/spec"
	"github.com/jpleatherland/spacetraders/internal/transforms"
)

func (s Server) GetContracts(w http.ResponseWriter, r *http.Request, params spec.GetContractsParams) {
	resources, ok := middleware.GetResources(r.Context())
	if !ok {
		response.RespondWithError(w, "couldn't get resources", http.StatusInternalServerError)
		return
	}

	session, ok := middleware.GetSession(r.Context())
	if !ok {
		response.RespondWithError(w, "couldn't get session", http.StatusInternalServerError)
		return
	}

	contractsResponse, ok := resources.Cache.Get("myContracts")
	if ok {
		response.RespondWithTemplate(w, "contracts.html", contractsResponse)
	}

	log.Println("not found in cache, getting contracts")
	//pageQuery := fmt.Sprintf("page=%v", params.Page)
	//limitQuery := fmt.Sprintf("limit=%v", params.Limit)
	url := baseUrl + "/my/contracts" // + pageQuery + "&" + limitQuery
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTML(w, "unable to get contracts", http.StatusInternalServerError)
		return
	}
	
	agentToken, err := resources.DB.GetAgentTokenById(r.Context(), session.AgentID.UUID)
	if err != nil {
		response.RespondWithError(w, "couldn't get agentid", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+agentToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTMLError(w, "unable to get contracts", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to unmarshal contract: ", err.Error())
		response.RespondWithHTMLError(w, "unable to get contract", http.StatusInternalServerError)
		return
	}

	contractsResponse = spec.ContractsResponse{}
	err = json.Unmarshal(body, &contractsResponse)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTMLError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, ok := contractsResponse.(spec.ContractsResponse)
	if !ok {
		log.Println("contracts response is not of type contracts response. possibly no contracts to get")
		response.RespondWithTemplate(w, "contracts.html", nil)
		return
	}

	resources.Cache.Add("myContracts", t, 0)
	response.RespondWithTemplate(w, "contracts.html", contractsResponse)
}

func (s Server) GetContract(w http.ResponseWriter, r *http.Request, contractId string) {
	resources, ok := middleware.GetResources(r.Context())
	if !ok {
		response.RespondWithError(w, "couldn't get resources", http.StatusInternalServerError)
		return
	}
	contractResponse, ok := resources.Cache.Get(contractId)
	if !ok {
		resp, err := http.Get("/my/contracts/" + contractId)
		if err != nil {
			log.Println("unable to get contract: ", err.Error())
			response.RespondWithHTMLError(w, "unable to get contract", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("unable to unmarshal contract: ", err.Error())
			response.RespondWithHTMLError(w, "unable to get contract", http.StatusInternalServerError)
			return
		}
		contractResponse = spec.ContractResponse{}
		err = json.Unmarshal(body, &contractResponse)
		resources.Cache.Add(contractId, contractResponse, 0)
	}
	t, ok := contractResponse.(spec.ContractResponse)
	if !ok {
		log.Println("contract response is not of type contract response")
		response.RespondWithHTMLError(w, "unable to get contract", http.StatusInternalServerError)
		return
	}
	htmlResp := transforms.StructureContract(t)
	response.RespondWithHTML(w, htmlResp, http.StatusOK)
}
func (s Server) AcceptContract(w http.ResponseWriter, r *http.Request, contractId string)  {}
func (s Server) DeliverContract(w http.ResponseWriter, r *http.Request, contractId string) {}
func (s Server) FulfillContract(w http.ResponseWriter, r *http.Request, contractId string) {}
