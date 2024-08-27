package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jpleatherland/spacetraders/internal/cache"
	"github.com/jpleatherland/spacetraders/internal/db"
)

type Resources struct {
	DB     *db.Queries
	Secret string
	Cache  *cache.Cache
}

type contextKey string

const (
	ResourcesKey = contextKey("resources")
	SessionKey   = contextKey("session")
)

func GetResources(ctx context.Context) (*Resources, bool) {
	res, ok := ctx.Value(ResourcesKey).(*Resources)
	return res, ok
}

func InjectResources(rsrc *Resources) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ResourcesKey, rsrc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SessionMiddleware is a middleware function for handling sessions
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example session handling
		// Here you might read/write session data, validate session tokens, etc.
		resources, ok := GetResources(r.Context())
		if !ok {
			http.Error(w, "Resources not found in context", http.StatusInternalServerError)
			return
		}

		sessionCookie, err := r.Cookie("spacetradersSession")
		if err != nil {
			log.Println("unable to find session cookie")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		session, err := resources.DB.GetSessionById(context.Background(), sessionCookie.Value)
		if err != nil {
			errMsg := fmt.Sprintf("unable to get session: %v", err.Error())
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), SessionKey, session)

		// Pass the request with the new context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetSessionID extracts the session ID from the context
func GetSession(ctx context.Context) (db.Session, bool) {
	session, ok := ctx.Value(SessionKey).(db.Session)
	return session, ok
}
