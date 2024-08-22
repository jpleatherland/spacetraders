package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/routes"
)

type contextKey string

func internalResourcesMiddleware(handler resourcesHandler, resources *routes.Resources) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req, resources)
	}
}

type resourcesHandler func(rw http.ResponseWriter, req *http.Request, resources *routes.Resources)

func sessionMiddleware(handler sessionHandler, resources *routes.Resources) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := context.Background()

		sessionCookie, err := req.Cookie("spacetradersSession")
		if err != nil {
			http.Redirect(rw, req, "/login", http.StatusFound)
			return
		}

		session, err := resources.DB.GetSessionById(ctx, sessionCookie.Value)
		if err != nil {
			errMsg := fmt.Sprintf("unable to get session: %v", err.Error())
			http.Error(rw, errMsg, http.StatusUnauthorized)
			return
		}

		handler(rw, req, session, resources)
	}
}

func redirectLogin(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		_, err := req.Cookie("spacetradersSession")
		if err == nil {
			http.Redirect(rw, req, "/home", http.StatusFound)
			return
		}
		handler(rw, req)
	}
}

type sessionHandler func(rw http.ResponseWriter, req *http.Request, session db.Session, resources *routes.Resources)

func resourcesMiddleware(resources interface{}) func(http.Handler) http.Handler {
	var resourceKey contextKey = "resources"
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), resourceKey, resources)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
