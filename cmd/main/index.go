package main

import (
	"log"
	"net/http"
)

func index(rw http.ResponseWriter, req *http.Request) {
	log.Println("I'm in index")
	isLoggedIn := checkIfUserLoggedIn(req)

	if isLoggedIn {
		http.Redirect(rw, req, "/home", 301)
	} else {
		http.Redirect(rw, req, "/login", 301)
	}
}

func checkIfUserLoggedIn(req *http.Request) bool {
	refreshToken := ""

	for _, cookie := range req.Cookies() {
		if cookie.Name == "spacetradersSession" {
			refreshToken = cookie.Value
		}
	}

	if refreshToken == "" {
		return false
	}

	return true
}
