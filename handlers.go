package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// main page handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

// check on health of server
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"alive": true}`)
}

// credentials struct to store user information
type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db"username"`
}

// register new users
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// there is something wrong with the json decode return error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// salt and hash the password using bcrypt. salt set to 8
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		// there is something wrong with the password hashing
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// insert new user into db
	if _, err = db.Query("INSERT INTO users (name, hash) VALUES ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
		// there is an error entering into db return 500
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// decode and store post request json
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// search db for user
	rows := db.QueryRow("SELECT hash FROM users WHERE name=$1", creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get stored password and compare with entered password
	storedCreds := &Credentials{}
	err = rows.Scan(&storedCreds.Password)
	if err != nil {
		// if no entry exists then deny login
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// other errors return 500
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// compare stored password with hash
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
