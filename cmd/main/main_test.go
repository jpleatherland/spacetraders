package main

import (
	"testing"
	"os"

	"github.com/jpleatherland/spacetraders/internal/db"
)

var resources = Resources{}

func TestMain(m *testing.M) {
	database, jwt, err := setupTestEnvironment()
	if err != nil {
		panic(err)
	}
	dbQueries := db.New(database)
	resources.DB = dbQueries
	resources.Secret = jwt
	code := m.Run()
	database.Close()
	teardownTestEnvironment()
	os.Exit(code)
}
