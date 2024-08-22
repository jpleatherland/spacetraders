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

	resources := Resources{
		DB:     dbQueries,
		Secret: JWT,
		Cache:  cache,
	}

	s := api.NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /css/output.css", func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, "./views/css/output.css")
	})

	mux.HandleFunc("POST /createUser", internalResourcesMiddleware(routes.CreateUser, &resources))
	mux.HandleFunc("POST /userlogin", internalResourcesMiddleware(routes.UserLogin, &resources))
	mux.HandleFunc("GET /home", sessionMiddleware(web.HomePage, &resources))
	mux.HandleFunc("GET /login", redirectLogin(web.LoginPage))
	mux.HandleFunc("GET /", index)

	resourcesHandler := resourcesMiddleware(&resources)
	stMux := api.HandlerWithOptions(s, api.StdHTTPServerOptions{
		BaseRouter:  mux,
		Middlewares: []api.MiddlewareFunc{resourcesHandler},
	})

	resources.Server = s

	server := &http.Server{
		Addr:    ":8080",
		Handler: stMux,
	}

	log.Printf("Listening on %v", server.Addr)

	log.Fatal(server.ListenAndServe())
}
