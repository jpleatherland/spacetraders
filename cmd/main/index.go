package main

import (
	"log"
	"net/http"
)

func index(rw http.ResponseWriter, req *http.Request) {
	log.Println("I'm in index")
	isLoggedIn := checkIfUserLoggedIn(req)

	if isLoggedIn {
		http.Redirect(rw, req, "/home", http.StatusFound)
	} else {
		http.Redirect(rw, req, "/login", http.StatusFound)
	}
}

func checkIfUserLoggedIn(req *http.Request) bool {
	refreshToken := ""

	for _, cookie := range req.Cookies() {
		if cookie.Name == "spacetradersSession" {
			refreshToken = cookie.Value
		}
	}

	return refreshToken != ""
}
