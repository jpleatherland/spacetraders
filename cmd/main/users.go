package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func (resources *Resources) createUser(rw http.ResponseWriter, req *http.Request) {
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

	dbUser := db.CreateUserParams{
		ID: uuid.New(),
		Name: newUser.Name,
		Password: string(hashedPassword[:]),
	}
	ctx := context.Background()

	_, err = resources.DB.CreateUser(ctx, dbUser)
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

func (resources *Resources) userLogin(rw http.ResponseWriter, req *http.Request) {
	userLogin, err := FormUserToStruct(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	user, err := resources.DB.GetUserByName(ctx, userLogin.Name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		http.Error(rw, "invalid password", http.StatusUnauthorized)
		return
	}

	token, expiryEpoch, err := generateToken(user.Name, 0, resources.Secret)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	//err = resources.DB.WriteRefreshToken(token, user.Token, expiryEpoch)
	//if err != nil {
	//   http.Error(rw, err.Error(), http.StatusInternalServerError)
	//

	rw.WriteHeader(http.StatusCreated)
	refreshTokenCookie := createRefreshTokenCookie(token, expiryEpoch)
	http.SetCookie(rw, &refreshTokenCookie)
	http.Redirect(rw, req, "/homepage", 301)
}

func UserToStruct(req *http.Request) (db.User, error) {
	var newUser db.User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func FormUserToStruct(req *http.Request) (db.User, error) {
	var newUser db.User
	err := req.ParseForm()
	if err != nil {
		return newUser, err
	}
	newUser.Name = req.FormValue("username")
	newUser.Password = req.FormValue("password")
	return newUser, nil
}
func createRefreshTokenCookie(token string, expiryTime int64) http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = "refreshToken"
	cookie.Value = token
	cookie.Expires = time.Unix(expiryTime, 0)
	return cookie
}
