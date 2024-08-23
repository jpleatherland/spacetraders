package routes

import (
	"testing"
	"os"

	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/middleware"

)

var resources = middleware.Resources{}

func TestMain(m *testing.M) {
	database, jwt, err := setupTestEnvironment()
	if err != nil {
		panic(err)
	}
	dbQueries := db.New(database)
	resources := middleware.Resources{}
	resources.DB = dbQueries
	resources.Secret = jwt
	code := m.Run()
	database.Close()
	teardownTestEnvironment()
	os.Exit(code)
}
