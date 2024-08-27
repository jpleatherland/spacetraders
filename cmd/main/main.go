package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/jpleatherland/spacetraders/internal/api"
	"github.com/jpleatherland/spacetraders/internal/cache"
	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/middleware"
	"github.com/jpleatherland/spacetraders/internal/routes"
	"github.com/jpleatherland/spacetraders/internal/spec"
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

	resources := middleware.Resources{
		DB:     dbQueries,
		Secret: JWT,
		Cache:  cache,
	}

	s := api.NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /stylesheet.css", func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, "./views/css/stylesheet.css")
	})
	mux.HandleFunc("GET /favicon.ico", func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, "./views/rocket-16.ico")
	})

	homeHandler := middleware.InjectResources(&resources)(
		middleware.SessionMiddleware(
			http.HandlerFunc(web.HomePage),
		),
	)

	mux.HandleFunc("POST /createUser", routes.CreateUser)
	mux.HandleFunc("POST /userlogin", routes.UserLogin)
	mux.Handle("GET /home", homeHandler)
	mux.HandleFunc("GET /login", redirectLogin(web.LoginPage))
	mux.HandleFunc("GET /", index)

	stMux := spec.HandlerWithOptions(s, spec.StdHTTPServerOptions{
		BaseRouter:  mux,
		Middlewares: []spec.MiddlewareFunc{},
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "login") || r.URL.Path == "/createUser" {
			middleware.InjectResources(&resources)(stMux).ServeHTTP(w, r)
			return
		}
		middleware.InjectResources(&resources)(
			middleware.SessionMiddleware(
				stMux),
		).ServeHTTP(w, r)
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Printf("Listening on %v", server.Addr)

	log.Fatal(server.ListenAndServe())
}

func redirectLogin(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Println("redirect login")
		_, err := req.Cookie("spacetradersSession")
		if err == nil {
			http.Redirect(rw, req, "/home", http.StatusFound)
			return
		}
		handler(rw, req)
	}
}
