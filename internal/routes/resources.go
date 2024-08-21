package routes

import (
    "github.com/jpleatherland/spacetraders/internal/db"
    "github.com/jpleatherland/spacetraders/internal/cache"
)

type Resources struct {
	DB     *db.Queries
	Secret string
	Cache  *cache.Cache
}
