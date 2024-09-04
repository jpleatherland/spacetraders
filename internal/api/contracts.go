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
