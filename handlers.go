package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/oliverpauffley/chess_ladder/models"
	_ "github.com/oliverpauffley/chess_ladder/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

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
	authRouter.HandleFunc("/users/{id:?[0-9]+}", env.UserHandler).Methods("GET")

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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// compare stored password with hash
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Hash), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// user is ok so authenticate
	user := User{ID: storedCreds.Id, Authenticated: true}
	session.Values["user"] = user
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
	session.Values["user"] = User{Authenticated: false}
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (env Env) UserHandler(w http.ResponseWriter, r *http.Request) {
	// create empty user credentials to return on an error
	emptyCredentials := models.CredentialsExternal{}
	emptyUser, _ := json.Marshal(emptyCredentials)

	// set http header type as json
	w.Header().Set("Content-Type", "application/json")

	// get user id from url entered
	vars := mux.Vars(r)
	log.Print(vars)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(emptyUser)
		return
	}

	// get user credentials from db and check for errors
	credentials, err := env.db.QueryById(id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(emptyUser)
		return
	}
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(emptyUser)
		return
	}

	// write credentials to json and return
	js, _ := json.Marshal(credentials)
	_, err = w.Write(js)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// the front end should send the following to login and register
type jsonCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

// User stores session information for a cookie
type User struct {
	ID            int  `json:"ID"`
	Authenticated bool `json:"authenticated"`
}

// GetUser is a helper function to return a user from the session value.
// if no user is found an empty unauthenticated user is returned
func GetUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}
