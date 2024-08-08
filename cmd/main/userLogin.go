package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	newUser, err := UserToStruct(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 4)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := createAgent(newUser.Username, newUser.Faction)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	dbUser := db.User{
		Username: newUser.Username,
		Password: hashedPassword,
		Faction:  newUser.Faction,
		Token:    token,
	}

	err = apiConf.DB.CreateUser(dbUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusConflict)
		return
	}
	response, err := json.Marshal(struct {
		Response string `json:"response"`
	}{
		Response: "user created successfully. please login to continue",
	})
	if err != nil {
		log.Println("unable to marshal response")
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write(response)
}

func (apiConf *apiConfig) userLogin(rw http.ResponseWriter, req *http.Request) {
	userLogin, err := UserToStruct(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := apiConf.DB.ReadUserFromDB(userLogin.Username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}
	fmt.Println(user)
	rw.WriteHeader(http.StatusCreated)
}

func GenerateToken(username string, tokenExpiryTime int32, jwtSecret string) (string, error) {
	expiryTime := 24 * time.Hour
	if tokenExpiryTime > 0 {
		expiryTime = time.Duration(tokenExpiryTime) * time.Second
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "jpleatherland/spacetraders",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryTime)),
		Subject:   username,
	})
	signedToken, err := newToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func UserToStruct(req *http.Request) (User, error) {
	var newUser User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}
