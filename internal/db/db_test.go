package db

import (
	"os"
	"database/sql"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/pressly/goose"
)

func setupTestEnvironment() (*sql.DB, *database.Queries, error) {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	DB_CONN_STRING := os.Getenv("DB_CONN_TEST")
	pgdb, err := sql.Open("postgres", DB_CONN_STRING+"?sslmode=disable")
	if err != nil {
		return err
	}
	defer pgdb.Close()
	pgdb.Exec("CREATE DATABASE spacetraders_test")
	db, err := sql.Open("postgres", DB_CONN_STRING+"/spacetraders_test?sslmode=disable")
	if err != nil {
		return err
	}
	err = goose.Up(db, "./sql/schema")
	if err != nil {
		return err
	}

	return nil
}

func teardownTestEnvironment(resources TestResources) {
	if resources.DBLocation != "" {
		os.Remove(resources.DBLocation)
	}
}

func setupDB() error {

	return nil
}

func TestCreateUser(t *testing.T) {
	resources, err := setupTestEnvironment()
	if err != nil {
		t.Errorf("failed to set up test environment: %v", err)
	}

	defer teardownTestEnvironment(*resources)

	_, err = CreateUser(*resources)
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}
}

func TestReadUser(t *testing.T) {
	resources, err := setupTestEnvironment()
	if err != nil {
		t.Errorf("failed to set up test environment: %v", err)
	}
	defer teardownTestEnvironment(*resources)
	user, err := CreateUser(*resources)
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	dbUser, err := ReadUser(user.Username, *resources)
	if err != nil {
		t.Errorf("failed to read from db: %v", err)
	}
	if dbUser.Username != user.Username {
		t.Errorf("incorrect username: got %v, want %v", dbUser.Username, user.Username)
	}
}

func TestWriteRefreshToken(t *testing.T) {
	resources, err := setupTestEnvironment()
	if err != nil {
		t.Errorf("failed to set up test environment: %v", err)
	}
	defer teardownTestEnvironment(*resources)

	err = resources.DB.WriteRefreshToken("abcdefg", "efghijk", time.Now().Add(time.Duration(24 * time.Hour)).Unix())
	if err != nil {
		t.Errorf("failed to write refresh token to db: %v",  err)
	}
}

func TestReadRefreshToken(t *testing.T) {
	resources, err := setupTestEnvironment()
	if err != nil {
		t.Errorf("failed to set up test environment: %v", err)
	}
	defer teardownTestEnvironment(*resources)

	refreshToken := "abcdefg"
	accessToken := "efghijk"

	err = resources.DB.WriteRefreshToken(refreshToken, accessToken, time.Now().Add(time.Duration(24 * time.Hour)).Unix())
	if err != nil {
		t.Errorf("failed to write refresh token to db: %v",  err)
	}

	dbToken, err := resources.DB.ReadRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("refresh token read failed: %v ", err)
	}
	if dbToken != accessToken {
		t.Errorf("got %v, want %v", dbToken, accessToken)
	}

}

func CreateUser(resources TestResources) (User, error) {
	newUser := User{
		Username: "User1234",
		Password: []byte("password"),
		Faction:  "COSMIC",
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGlmaWVyIjoiVEVTVF8yNzUzIiwidmVyc2lvbiI6InYyLjIuMCIsInJlc2V0X2RhdGUiOiIyMDI0LTA3LTIxIiwiaWF0IjoxNzIyODY4OTAwLCJzdWIiOiJhZ2VudC10b2tlbiJ9.mKGCgpYsRNezk1H_7tke0Vf0hv-aaTYN4N8HJpeULDlDstKHXX11iUgBGpb7exzejRzQom6czQSdpz-Sr1tvggcISujPnwyRmbtDOrf1lEPsvfbsGNyNFZvpXq2vaDUQY8Incm_6-1cxekRHmH76KEdq1mpPRU_GseAVxqHBnMG87W6HbFpb-cksGr8mlb8XvPcQruFXbNwexItkTJKtT9kxv35-54X_RwZ_9efcE3mR5cu1wvf3SRJ4E1mn5VB6jW-Gf82m--QM2Znj-Lb58x_i-DDkUK3kJPP8RgiaO6Mu0NPYdsNDojnA8mc2tcr0_gPqEwnnkHTez0l2uSQOnQ",
	}

	err := resources.DB.CreateUser(newUser)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func ReadUser(username string, resources TestResources) (User, error) {
	user, err := resources.DB.ReadUserFromDB(username)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
