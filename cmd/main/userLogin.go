package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Username string `json:"symbol"`
	Password string `json:"password"`
	Faction  string `json:"faction"`
	Token    string `json:"token"`
}

func userLogin(rw http.ResponseWriter, req *http.Request) {
}

func (apiConf *apiConfig) createUser(rw http.ResponseWriter, req *http.Request) {
	var newUser User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}
