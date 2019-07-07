package main

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/oliverpauffley/chess_ladder/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	authRouter.HandleFunc("/users/{id:[0-9]+}", env.UserHandler).Methods("GET")

	return router
}

// register new users
func (env Env) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	register := &RegisterCredentials{}
	err := json.NewDecoder(r.Body).Decode(register)
	if err != nil {
		// there is something wrong with the json decode return error
		log.Printf("There was an error with json decoding, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// email must be lowercase
	register.Email = strings.ToLower(register.Email)

	// validate the json credentials
	if register.Username == "" || register.Password == "" || register.Password != register.Confirm {
		log.Print(register)
		log.Printf("There was a problem validating the json creds, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//check if user currently exists and check if email is already in db
	if _, err := env.db.QueryByName(register.Username); err != sql.ErrNoRows {
		log.Print("Trying to register a user that already exists")
		w.WriteHeader(http.StatusConflict)
		return
	}

	// use create user to send request to server
	err = env.db.CreateUser(register.Username, register.Email, register.Password)
	if err != nil {
		log.Printf("There was an error creating a user, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Login existing users and provide them with a cookie
func (env Env) LoginHandler(w http.ResponseWriter, r *http.Request) {

	// decode and store post request json
	creds := &LoginCredentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// search db for user
	storedCreds, err := env.db.QueryByName(creds.Username)
	if err == sql.ErrNoRows {
		log.Printf("name: %s", creds.Username)
		log.Print("No User exists with this name")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// compare stored password with hash
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Hash), []byte(creds.Password)); err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expireTime := time.Now().Add(24 * time.Hour)
	// user is ok so authenticate
	user := User{
		storedCreds.Id,
		storedCreds.Username,
		jwt.StandardClaims{
			// token lasts 1 day
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ladderapp",
		},
	}

	// create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// set user cookie to token value
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expireTime,
		HttpOnly: false,
		Secure:   true,
	})
	_, _ = w.Write([]byte(tokenString))
}

// Logout a logged in user
func (env Env) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   " ",
		Expires: time.Now(),
	})
}

func (env Env) UserHandler(w http.ResponseWriter, r *http.Request) {
	// set http header type as json
	w.Header().Set("Content-Type", "application/json")

	// get user id from url entered
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get user credentials from db and check for errors
	credentials, err := env.db.QueryById(id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write credentials to json and return
	js, _ := json.Marshal(credentials)
	_, err = w.Write(js)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// the front end should send the following to login and register
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterCredentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

// Custom jwt claims struct that inherits from standard
type User struct {
	ID       int    `json:"ID"`
	Username string `json:"username"`
	jwt.StandardClaims
}
