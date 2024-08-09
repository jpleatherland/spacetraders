package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	newUser, err := FormUserToStruct(req)
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
	userLogin, err := FormUserToStruct(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := apiConf.DB.ReadUserFromDB(userLogin.Username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(userLogin.Password))
	if err != nil {
		http.Error(rw, "invalid password", http.StatusUnauthorized)
		return
	}

	token, expiryEpoch, err := generateToken(user.Username, 0, apiConf.Secret)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	err = apiConf.DB.WriteRefreshToken(token, user.Token, expiryEpoch)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusCreated)
	refreshTokenCookie := createRefreshTokenCookie(token, expiryEpoch)
	http.SetCookie(rw, &refreshTokenCookie)
	http.Redirect(rw, req, "/homepage", 301)
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

func FormUserToStruct(req *http.Request) (User, error) {
	var newUser User
	err := req.ParseForm()
	if err != nil {
		return newUser, err
	}
	newUser.Username = req.FormValue("uname")
	newUser.Password = req.FormValue("psw")
	newUser.Faction = req.FormValue("faction")
	return newUser, nil
}
func createRefreshTokenCookie(token string, expiryTime int64) http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = "refreshToken"
	cookie.Value = token
	cookie.Expires = time.Unix(expiryTime, 0)
	return cookie
}
