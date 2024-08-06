package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jpleatherland/spacetraders/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"symbol"`
	Password string `json:"password"`
	Faction  string `json:"faction"`
	Token    string `json:"token"`
}

func (apiConf *apiConfig) createUser(rw http.ResponseWriter, req *http.Request) {
	var newUser User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	newToken, err := GenerateToken(newUser.Username, 100, apiConf.Secret)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 4)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	createAgent(newUser.Username, newUser.Faction)

	dbUser := db.User{
		Username: newUser.Username,
		Password: hashedPassword,
		Faction:  newUser.Faction,
		Token:    newToken,
	}

	err = apiConf.DB.CreateUser(dbUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusConflict)
	}
	
	rw.WriteHeader(http.StatusCreated)
}

func (apiConf *apiConfig) userLogin(rw http.ResponseWriter, req *http.Request) {
}

func GenerateToken(userEmail string, tokenExpiryTime int32, jwtSecret string) (string, error) {
	expiryTime := 24 * time.Hour
	if tokenExpiryTime > 0 {
		expiryTime = time.Duration(tokenExpiryTime) * time.Second
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "jpleatherland/spacetraders",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryTime)),
		Subject:   userEmail,
	})
	signedToken, err := newToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
