package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/response"
	"github.com/jpleatherland/spacetraders/internal/spec"
)

func (s Server) GetMyShips(w http.ResponseWriter, r *http.Request, params spec.GetMyShipsParams) {
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

	cachedShipsResponse, ok := resources.Cache.Get("myShips")
	if ok {
		response.RespondWithPartialTemplate(w, "fleetPartials", "myships.html", cachedShipsResponse, nil)
		return
	}

	log.Println("not found in cache, getting ships")
	//pageQuery := fmt.Sprintf("page=%v", params.Page)
	//limitQuery := fmt.Sprintf("limit=%v", params.Limit)
	url := baseUrl + "/my/ships" // + pageQuery + "&" + limitQuery
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTML(w, "unable to get ships", http.StatusInternalServerError)
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
		response.RespondWithHTMLError(w, "unable to get ships", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to unmarshal ships: ", err.Error())
		response.RespondWithHTMLError(w, "unable to get ships", http.StatusInternalServerError)
		return
	}

	shipsResponse := spec.ShipsResponse{}
	err = json.Unmarshal(body, &shipsResponse)
	if err != nil {
		log.Println(err.Error())
		response.RespondWithHTMLError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resources.Cache.Add("myShips", shipsResponse, 0)
	response.RespondWithPartialTemplate(w, "fleetPartials", "myships.html", shipsResponse, nil)
}
