package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http/httptest"
	"net/http"
	"os"
	"bytes"
	"fmt"
	"errors"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func setupTestEnvironment() (*sql.DB, string, error) {
	godotenv.Load("../../.env")
	DB_CONN_STRING := os.Getenv("DB_CONN_TEST")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	ROOT_DIR := os.Getenv("ROOT_DIR")
	pgdb, err := sql.Open("postgres", DB_CONN_STRING)
	if err != nil {
		return nil, "", err
	}

	_, err = pgdb.Exec("CREATE DATABASE spacetraders_test")
	if err != nil {
		_, err = pgdb.Exec("DROP DATABASE spacetraders_test")
		if err != nil {
			return nil, "", err
		}
		log.Println("dropped test table")
		_, err = pgdb.Exec("CREATE DATABASE spacetraders_test")
		if err != nil {
			return nil, "", err
		}
	}

	testdb, err := sql.Open("postgres", DB_CONN_STRING+" database=spacetraders_test")
	if err != nil {
		return nil, "", err
	}

	err = goose.Up(testdb, ROOT_DIR+"/sql/schema")
	if err != nil {
		return nil, "", err
	}

	return testdb, JWT_SECRET, nil
}

func teardownTestEnvironment() {
	godotenv.Load()
	DB_CONN_STRING := os.Getenv("DB_CONN_TEST")
	pgdb, err := sql.Open("postgres", DB_CONN_STRING)
	if err != nil {
		log.Printf("failed to reopen database: %v", err.Error())
	}

	_, err = pgdb.Exec("DROP DATABASE spacetraders_test")
	if err != nil {
		log.Printf("failed to teardown test environment: %v", err.Error())
	}
	log.Println("dropped test database")
}

func createTestUser(resources *Resources, name, password string) error {
	baseURL := "http://localhost:8080"
	userBody := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}
	payload, err := json.Marshal(userBody)
	if err != nil {
		return err
	}

	req := httptest.NewRequest("POST", baseURL, bytes.NewBuffer(payload))
	res := httptest.NewRecorder()
	resources.createUser(res, req)

	if res.Result().StatusCode != http.StatusCreated{
		errorMsg := fmt.Sprintf("create user response status incorrect, expected: %v got: %v", http.StatusCreated, res.Result().Status)
		return errors.New(errorMsg)
	}
	return nil
}
