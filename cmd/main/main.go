package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jpleatherland/spacetraders/internal/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	jwt := os.Getenv("JWT_SECRET")
	dbConnection, err := db.NewDB(os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	apiConf := apiConfig{
		DB:     dbConnection,
		Secret: jwt,
	}
	http.HandleFunc("POST /create", apiConf.createUser)
	http.HandleFunc("POST /login", apiConf.userLogin)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type apiConfig struct {
	DB     *db.DB
	Secret string
}
