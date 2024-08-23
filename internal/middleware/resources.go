package middleware

import (
    "context"
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
const resourcesKey = contextKey("resources")

func GetResources(ctx context.Context) (*Resources, bool) {
    res, ok := ctx.Value(resourcesKey).(*Resources)
    return res, ok
}

func InjectMiddleware(rsrc *Resources) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), resourcesKey, rsrc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
