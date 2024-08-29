package spec

type ServerStatus struct {
	ServerResets struct {
		Next      string `json:"next"`
		Frequency string `json:"frequency"`
	} `json:"serverResets"`
	RequestStatus string `json:"RequestStatus"`
	Status        string `json:"status"`
	Version       string `json:"version"`
	ResetDate     string `json:"resetDate"`
	Description   string `json:"description"`
	Leaderboards  struct {
		MostCredits []struct {
			AgentSymbol string `json:"agentSymbol"`
			Credits     int64  `json:"credits"`
		} `json:"mostCredits"`
		MostSubmittedCharts []struct {
			AgentSymbol string `json:"agentSymbol"`
			ChartCount  int    `json:"chartCount"`
		} `json:"mostSubmittedCharts"`
	} `json:"leaderboards"`
	Announcements []struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"announcements"`
	Links []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"links"`
	Stats struct {
		Agents    int `json:"agents"`
		Ships     int `json:"ships"`
		Systems   int `json:"systems"`
		Waypoints int `json:"waypoints"`
	} `json:"stats"`
}

type RegisterResponse struct {
	Data struct {
		Agent    Agent    `json:"agent"`
		Contract Contract `json:"contract"`
		Faction  Faction  `json:"faction"`
		Ship     Ship     `json:"ship"`
		Token    string   `json:"token"`
	} `json:"data"`
}

type AgentResponse struct {
	Data Agent `json:"data"`
}

var Factions = [19]FactionSymbol{
	FactionSymbolAEGIS,
	FactionSymbolANCIENTS,
	FactionSymbolASTRO,
	FactionSymbolCOBALT,
	FactionSymbolCORSAIRS,
	FactionSymbolCOSMIC,
	FactionSymbolCULT,
	FactionSymbolDOMINION,
	FactionSymbolECHO,
	FactionSymbolETHEREAL,
	FactionSymbolGALACTIC,
	FactionSymbolLORDS,
	FactionSymbolOBSIDIAN,
	FactionSymbolOMEGA,
	FactionSymbolQUANTUM,
	FactionSymbolSHADOW,
	FactionSymbolSOLITARY,
	FactionSymbolUNITED,
	FactionSymbolVOID,
}
