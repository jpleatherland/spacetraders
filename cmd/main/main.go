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
	DB_CONN_STRING := os.Getenv("DB_CONN")
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

	mux.HandleFunc("GET /output.css", func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, "./views/output.css")
		})

	mux.HandleFunc("POST /createUser", resources.createUser)
	mux.HandleFunc("POST /userlogin", resources.userLogin)
	mux.HandleFunc("GET /createForm", createForm)
	mux.HandleFunc("GET /login", resources.loginPage)
	mux.HandleFunc("GET /home", resources.sessionMiddleware(resources.homePage))
	mux.HandleFunc("GET /", resources.index)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Printf("Listening on %v", server.Addr)

	log.Fatal(server.ListenAndServe())
}

type Resources struct {
	DB     *db.Queries
	Secret string
}
