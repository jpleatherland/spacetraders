package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/jpleatherland/spacetraders/internal/db"
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
	}{
		Name:     name,
		Password: password,
	}
	payload, err := json.Marshal(userBody)
	if err != nil {
		return err
	}

	req := httptest.NewRequest("POST", baseURL, bytes.NewBuffer(payload))
	res := httptest.NewRecorder()
	resources.createUser(res, req)

	if res.Result().StatusCode != http.StatusCreated {
		errorMsg := fmt.Sprintf("create user response status incorrect, expected: %v got: %v", http.StatusCreated, res.Result().Status)
		return errors.New(errorMsg)
	}
	return nil
}

func testUserLogin(resources *Resources, name, password string) (*http.Cookie, error) {
	baseURL := "http://localhost:8080"
	data := url.Values{}
	data.Set("name", name)
	data.Set("password", password)

	req := httptest.NewRequest("POST", baseURL, strings.NewReader(data.Encode()))
	res := httptest.NewRecorder()
	resources.userLogin(res, req)

	expected := http.StatusCreated
	got := res.Result().StatusCode

	if got != expected {
		errorMsg := fmt.Sprintf("failed to login, expected: %v, got: %v", expected, got)
		return nil, errors.New(errorMsg)
	}

	for _, cookie := range res.Result().Cookies() {
		if cookie.Name == "spacetradersSession" {
			return cookie, nil
		}
	}

	return nil, errors.New("Session cookie not found")
}

func testGetSessionById(resources *Resources, cookie *http.Cookie) (db.Session, error) {
	ctx := context.Background()
	session := db.Session{}
	session, err := resources.DB.GetSessionById(ctx, cookie.Value)
	return session, err
}
