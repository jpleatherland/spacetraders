package main

import (
	"net/http"
)

func (resources *Resources) index(rw http.ResponseWriter, req *http.Request) {
	isLoggedIn := resources.checkIfUserLoggedIn(req)

	if isLoggedIn {
		http.Redirect(rw, req, "/homepage", 301)
	} else {
		http.Redirect(rw, req, "/login", 301)
	}
}

func (resources *Resources) checkIfUserLoggedIn(req *http.Request) bool {
	refreshToken := ""

	for _, cookie := range req.Cookies() {
		if cookie.Name == "refreshToken" {
			refreshToken = cookie.Value
		}
	}

	if refreshToken == "" {
		return false
	}

	//_, err := resources.DB.ReadRefreshToken(refreshToken)

	return false

}
