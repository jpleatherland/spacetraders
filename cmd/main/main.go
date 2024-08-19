package main

import (
	"log"
	"net/http"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/jpleatherland/spacetraders/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	JWT := os.Getenv("JWT_SECRET")
	DB_CONN_STRING :=os.Getenv("DB_CONN_STRING")
	database, err := sql.Open("postgres", DB_CONN_STRING)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := db.New(database)

	resources := Resources{
		DB:     dbQueries,
		Secret: JWT,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /createUser", resources.createUser)
	mux.HandleFunc("POST /userlogin", resources.userLogin)
	mux.HandleFunc("GET /createForm", createForm)
	mux.HandleFunc("GET /login", resources.loginPage)
	mux.HandleFunc("GET /home", resources.homePage)
	mux.HandleFunc("GET /", resources.index)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

type Resources struct {
	DB     *db.Queries
	Secret string
}
