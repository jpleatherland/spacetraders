package routes

import (
	"github.com/jpleatherland/spacetraders/internal/api"
	"github.com/jpleatherland/spacetraders/internal/cache"
	"github.com/jpleatherland/spacetraders/internal/db"
)

type Resources struct {
	DB     *db.Queries
	Secret string
	Cache  *cache.Cache
	Server api.Server
}
