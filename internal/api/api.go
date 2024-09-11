package api

import (
	"log"
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/spec"
)

var baseUrl string = "https://api.spacetraders.io/v2"

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) GetFactions(w http.ResponseWriter, r *http.Request, params spec.GetFactionsParams) {}
func (s Server) GetFaction(w http.ResponseWriter, r *http.Request, factionSymbol string)           {}
func (s Server) PurchaseShip(w http.ResponseWriter, r *http.Request)                             {}
func (s Server) GetMyShip(w http.ResponseWriter, r *http.Request, shipSymbol string)             {}
func (s Server) GetMyShipCargo(w http.ResponseWriter, r *http.Request, shipSymbol string)        {}
func (s Server) CreateChart(w http.ResponseWriter, r *http.Request, shipSymbol string)           {}
func (s Server) GetShipCooldown(w http.ResponseWriter, r *http.Request, shipSymbol string)       {}
func (s Server) DockShip(w http.ResponseWriter, r *http.Request, shipSymbol string)              {}
func (s Server) ExtractResources(w http.ResponseWriter, r *http.Request, shipSymbol string)      {}
func (s Server) ExtractResourcesWithSurvey(w http.ResponseWriter, r *http.Request, shipSymbol string) {
}
func (s Server) Jettison(w http.ResponseWriter, r *http.Request, shipSymbol string)               {}
func (s Server) JumpShip(w http.ResponseWriter, r *http.Request, shipSymbol string)               {}
func (s Server) GetMounts(w http.ResponseWriter, r *http.Request, shipSymbol string)              {}
func (s Server) InstallMount(w http.ResponseWriter, r *http.Request, shipSymbol string)           {}
func (s Server) RemoveMount(w http.ResponseWriter, r *http.Request, shipSymbol string)            {}
func (s Server) GetShipNav(w http.ResponseWriter, r *http.Request, shipSymbol string)             {}
func (s Server) PatchShipNav(w http.ResponseWriter, r *http.Request, shipSymbol string)           {}
func (s Server) NavigateShip(w http.ResponseWriter, r *http.Request, shipSymbol string)           {}
func (s Server) NegotiateContract(w http.ResponseWriter, r *http.Request, shipSymbol string)      {}
func (s Server) OrbitShip(w http.ResponseWriter, r *http.Request, shipSymbol string)              {}
func (s Server) PurchaseCargo(w http.ResponseWriter, r *http.Request, shipSymbol string)          {}
func (s Server) ShipRefine(w http.ResponseWriter, r *http.Request, shipSymbol string)             {}
func (s Server) RefuelShip(w http.ResponseWriter, r *http.Request, shipSymbol string)             {}
func (s Server) GetRepairShip(w http.ResponseWriter, r *http.Request, shipSymbol string)          {}
func (s Server) RepairShip(w http.ResponseWriter, r *http.Request, shipSymbol string)             {}
func (s Server) CreateShipShipScan(w http.ResponseWriter, r *http.Request, shipSymbol string)     {}
func (s Server) CreateShipSystemScan(w http.ResponseWriter, r *http.Request, shipSymbol string)   {}
func (s Server) CreateShipWaypointScan(w http.ResponseWriter, r *http.Request, shipSymbol string) {}
func (s Server) GetScrapShip(w http.ResponseWriter, r *http.Request, shipSymbol string)           {}
func (s Server) ScrapShip(w http.ResponseWriter, r *http.Request, shipSymbol string)              {}
func (s Server) SellCargo(w http.ResponseWriter, r *http.Request, shipSymbol string)              {}
func (s Server) SiphonResources(w http.ResponseWriter, r *http.Request, shipSymbol string)        {}
func (s Server) CreateSurvey(w http.ResponseWriter, r *http.Request, shipSymbol string)           {}
func (s Server) TransferCargo(w http.ResponseWriter, r *http.Request, shipSymbol string)          {}
func (s Server) WarpShip(w http.ResponseWriter, r *http.Request, shipSymbol string)               {}
func (s Server) GetSystems(w http.ResponseWriter, r *http.Request, params spec.GetSystemsParams)  {}
func (s Server) GetSystem(w http.ResponseWriter, r *http.Request, systemSymbol string)            {}
func (s Server) GetSystemWaypoints(w http.ResponseWriter, r *http.Request, systemSymbol string, params spec.GetSystemWaypointsParams) {
}
func (s Server) GetWaypoint(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
	log.Println("in get waypoint")
	log.Println(systemSymbol, waypointSymbol)
}
func (s Server) GetConstruction(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
}
func (s Server) SupplyConstruction(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
}
func (s Server) GetJumpGate(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
}
func (s Server) GetMarket(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
}
func (s Server) GetShipyard(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
}
