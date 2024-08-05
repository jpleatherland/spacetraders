package main

import (
	"log"
	"net/http"

	"github.com/jpleatherland/internal/db"
)

func main() {
	dbConnection, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	apiConf := apiConfig{
		DB: dbConnection,
	}
	http.HandleFunc("POST /create", apiConf.createUser)
	http.HandleFunc("POST /login", userLogin)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type apiConfig struct {
	DB *db.DB
}
