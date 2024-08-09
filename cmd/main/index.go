package main

import (
	"net/http"
)

func (apiConf *apiConfig) index(rw http.ResponseWriter, req *http.Request) {
	isLoggedIn := apiConf.checkIfUserLoggedIn(req)

	if isLoggedIn {
		http.Redirect(rw, req, "/homepage", 301)
	} else {
		http.Redirect(rw, req, "/login", 301)
	}
}

func (apiConf *apiConfig) checkIfUserLoggedIn(req *http.Request) bool {
	refreshToken := ""

	for _, cookie := range req.Cookies() {
		if cookie.Name == "refreshToken" {
			refreshToken = cookie.Value
		}
	}

	if refreshToken == "" {
		return false
	}

	_, err := apiConf.DB.ReadRefreshToken(refreshToken)

	return err == nil

}
