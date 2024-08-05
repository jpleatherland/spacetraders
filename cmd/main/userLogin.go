package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jpleatherland/spacetraders/internal/db"
)

type User struct {
	Username string `json:"symbol"`
	Password string `json:"password"`
	Faction  string `json:"faction"`
	Token    string `json:"token"`
}

func (apiConf *apiConfig) userLogin(rw http.ResponseWriter, req *http.Request) {
}

func (apiConf *apiConfig) createUser(rw http.ResponseWriter, req *http.Request) {
	var newUser User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	newToken, err := GenerateToken(newUser.Username, 100, "lsakjdfoiasdf")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	dbUser := db.User{
		Username: newUser.Username,
		Password: newUser.Password,
		Faction:  newUser.Faction,
		Token:    newToken,
	}
	apiConf.DB.CreateUser(dbUser)

}

func GenerateToken(userEmail string, tokenExpiryTime int32, jwtSecret string) (string, error) {
	expiryTime := 24 * time.Hour
	if tokenExpiryTime > 0 {
		expiryTime = time.Duration(tokenExpiryTime) * time.Second
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
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
