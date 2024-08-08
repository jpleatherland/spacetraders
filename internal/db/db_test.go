package db

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/joho/godotenv"
)

type TestResources struct {
	DB         *DB
	DBLocation string
	Secret     string
}

func setupTestEnvironment() (*TestResources, error) {
	var tr TestResources
	// create the tempdb at path in .env
	err := godotenv.Load()
	if err != nil {
		return &tr, err
	}
	err = setupDB()
	if err != nil {
		return &tr, err
	}
	// connect to db
	db, err := NewDB(os.Getenv("DB_PATH"))
	if err != nil {
		return &tr, err
	}
	tr.DB = db
	tr.DBLocation = path.Join(os.Getenv("DB_PATH"), "spacetraders.db")

	return &tr, nil
}

func teardownTestEnvironment(resources TestResources) {
	if resources.DBLocation != "" {
		os.Remove(resources.DBLocation)
	}
}

func setupDB() error {
	cmd := exec.Command("bash", "-c", "./create_test_db.sh")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
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
