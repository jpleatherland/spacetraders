package api

import (
	"net/http"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) GetFactions(w http.ResponseWriter, r *http.Request, params GetFactionsParams)   {}
func (s Server) GetFaction(w http.ResponseWriter, r *http.Request, factionSymbol string)        {}
func (s Server) GetContracts(w http.ResponseWriter, r *http.Request, params GetContractsParams) {}
func (s Server) GetContract(w http.ResponseWriter, r *http.Request, contractId string)          {}
func (s Server) AcceptContract(w http.ResponseWriter, r *http.Request, contractId string)       {}
func (s Server) DeliverContract(w http.ResponseWriter, r *http.Request, contractId string)      {}
func (s Server) FulfillContract(w http.ResponseWriter, r *http.Request, contractId string)      {}
func (s Server) GetMyShips(w http.ResponseWriter, r *http.Request, params GetMyShipsParams)     {}
func (s Server) PurchaseShip(w http.ResponseWriter, r *http.Request)                            {}
func (s Server) GetMyShip(w http.ResponseWriter, r *http.Request, shipSymbol string)            {}
func (s Server) GetMyShipCargo(w http.ResponseWriter, r *http.Request, shipSymbol string)       {}
func (s Server) CreateChart(w http.ResponseWriter, r *http.Request, shipSymbol string)          {}
func (s Server) GetShipCooldown(w http.ResponseWriter, r *http.Request, shipSymbol string)      {}
func (s Server) DockShip(w http.ResponseWriter, r *http.Request, shipSymbol string)             {}
func (s Server) ExtractResources(w http.ResponseWriter, r *http.Request, shipSymbol string)     {}
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
func (s Server) Register(w http.ResponseWriter, r *http.Request)                                  {}
func (s Server) GetSystems(w http.ResponseWriter, r *http.Request, params GetSystemsParams)       {}
func (s Server) GetSystem(w http.ResponseWriter, r *http.Request, systemSymbol string)            {}
func (s Server) GetSystemWaypoints(w http.ResponseWriter, r *http.Request, systemSymbol string, params GetSystemWaypointsParams) {
}
func (s Server) GetWaypoint(w http.ResponseWriter, r *http.Request, systemSymbol string, waypointSymbol string) {
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
