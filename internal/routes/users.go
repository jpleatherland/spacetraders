package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	
	"github.com/google/uuid"
	"github.com/jpleatherland/spacetraders/internal/db"
	"github.com/jpleatherland/spacetraders/internal/response"
)

func CreateUser(rw http.ResponseWriter, req *http.Request, resources *Resources) {
	newUser, err := formUserToStruct(req)
	if err != nil {
		response.RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if newUser.Name == "" || newUser.Password == "" {
		response.RespondWithHTMLError(rw, "Please enter username and password", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 4)
	if err != nil {
		response.RespondWithHTMLError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	dbUser := db.CreateUserParams{
		ID:       uuid.New(),
		Name:     newUser.Name,
		Password: string(hashedPassword[:]),
	}
	ctx := context.Background()

	_, err = resources.DB.CreateUser(ctx, dbUser)
	if err != nil {
		response.RespondWithHTMLError(rw, err.Error(), http.StatusConflict)
		return
	}
	response.RespondWithHTML(rw, "<p>User created successfully, please login to continue</p>", http.StatusCreated)
	//response.RespondWithJSON(rw, http.StatusCreated, map[string]string{"Response": "user created successfully, please login to continue"})
}

func UserLogin(rw http.ResponseWriter, req *http.Request, resources *Resources) {
	userLogin, err := formUserToStruct(req)
	if err != nil {
		response.RespondWithError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	user, err := resources.DB.GetUserByName(ctx, userLogin.Name)
	if err != nil {
		response.RespondWithError(rw, err.Error(), http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		response.RespondWithError(rw, "invalid password", http.StatusUnauthorized)
		return
	}

	token, expiryEpoch, err := generateToken(user.Name, 0, resources.Secret)
	if err != nil {
		response.RespondWithError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionCookie := createSessionCookie(token, expiryEpoch)
	http.SetCookie(rw, &sessionCookie)

	session := db.CreateSessionParams{
		ID:        token,
		ExpiresAt: sessionCookie.Expires,
		UserID:    user.ID,
		AgentID:   uuid.NullUUID{},
	}

	err = resources.DB.CreateSession(ctx, session)

	if err != nil {
		response.RespondWithError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("HX-Redirect", "http://localhost:8080/home")
	rw.WriteHeader(http.StatusOK)
}

func userToStruct(req *http.Request) (db.User, error) {
	var newUser db.User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func formUserToStruct(req *http.Request) (db.User, error) {
	var newUser db.User
	err := req.ParseForm()
	if err != nil {
		return newUser, err
	}
	newUser.Name = req.FormValue("username")
	newUser.Password = req.FormValue("password")
	return newUser, nil
}
