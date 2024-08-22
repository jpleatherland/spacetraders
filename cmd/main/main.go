package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jpleatherland/spacetraders/internal/api"
	"github.com/jpleatherland/spacetraders/internal/cache"
	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/web"
	_ "github.com/lib/pq"
)

const BASEURL = "api.spacetraders.io/v2"

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
	cache := cache.NewCache(5 * time.Minute)

	resources := routes.Resources{
		DB:     dbQueries,
		Secret: JWT,
		Cache:  cache,
	}

	s := api.NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /css/output.css", func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, "./views/css/output.css")
	})

	mux.HandleFunc("POST /createUser", resourcesMiddleware(routes.CreateUser, &resources))
	mux.HandleFunc("POST /userlogin", resourcesMiddleware(routes.UserLogin, &resources))
	mux.HandleFunc("GET /home", sessionMiddleware(web.HomePage, &resources))
	mux.HandleFunc("GET /login", redirectLogin(web.LoginPage))
	mux.HandleFunc("GET /", index)

	stMux := api.HandlerFromMux(s, mux)

	resources.Server = s

	server := &http.Server{
		Addr:    ":8080",
		Handler: stMux,
	}

	log.Printf("Listening on %v", server.Addr)

	log.Fatal(server.ListenAndServe())
}
