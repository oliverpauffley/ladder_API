package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/oliverpauffley/chess_ladder/models"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type jsonCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

// create a new router
func (env *Env) NewRouter() *mux.Router {
	router := mux.NewRouter()

	// set up routes
	router.HandleFunc("/register", env.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", env.LoginHandler).Methods("POST")

	// create authenticated routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Use(AuthMiddleware)
	authRouter.HandleFunc("/logout", env.LogoutHandler).Methods("GET")

	return router
}

// register new users
func (env Env) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	register := &jsonCredentials{}
	err := json.NewDecoder(r.Body).Decode(register)
	if err != nil {
		// there is something wrong with the json decode return error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate the json credentials
	if register.Username == "" || register.Password == "" || register.Password != register.Confirm {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//check if user currently exists
	if _, err := env.db.QueryByName(register.Username); err != sql.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// use create user to send request to server
	err = env.db.CreateUser(register.Username, register.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Login existing users and provide them with a cookie
func (env Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// create cookie ready to add to use
	session, err := store.Get(r, "authentication-cookie")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// decode and store post request json
	creds := &jsonCredentials{}
	err = json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// search db for user
	storedCreds, err := env.db.QueryByName(creds.Username)
	if err == sql.ErrNoRows {
		log.Print("No User exists with this name")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		fmt.Print(storedCreds)
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// compare stored password with hash
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Hash), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// user is ok so authenticate
	session.Values["authenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Logout a logged in user
func (env Env) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "authentication-cookie")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// revoke user cookie
	session.Values["authenticated"] = false
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
