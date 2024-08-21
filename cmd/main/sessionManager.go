package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/cache"
)

func createSessionCookie(token string, expiryTime int64) http.Cookie {
	cookie := http.Cookie{
		Name:       "spacetradersSession",
		Value:      token,
		Expires:    time.Unix(expiryTime, 0),
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
	}
	return cookie
}

func (resources *Resources) sessionMiddleware(handler sessionHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := context.Background()

		sessionCookie, err := req.Cookie("spacetradersSession")
		if err != nil {
			http.Error(rw, "no session cookie found", http.StatusUnauthorized)
			return
		}

		session, err := resources.DB.GetSessionById(ctx, sessionCookie.Value)
		if err != nil {
			errMsg := fmt.Sprintf("unable to get session: %v", err.Error())
			http.Error(rw, errMsg, http.StatusUnauthorized)
			return
		}

		handler(rw, req, session, resources.Cache)
	}
}

type sessionHandler func(rw http.ResponseWriter, req *http.Request, session db.Session, cache *cache.Cache)
