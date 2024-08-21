package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
