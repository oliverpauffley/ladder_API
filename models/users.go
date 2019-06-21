package models

import (
	. "golang.org/x/crypto/bcrypt"
	"time"
)

// credentials struct to store user information
type Credentials struct {
	Id       int       
	Hash     string    
	Username string    
	Wins     int       
	Draws    int       
	Losses   int       
	JoinDate time.Time 
}

func (db *DB) CreateUser(username, password string) error {
	// salt and hash the password using bcrypt. salt set to 8
	hashedPassword, err := GenerateFromPassword([]byte(password), 8)
	if err != nil {
		// there is something wrong with the password hashing
		return err
	}

	// insert new user into db
	_, err = db.Query("INSERT INTO users (name, hash) VALUES ($1, $2)",
		username, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) QueryByName(username string) (*Credentials, error) {
	rows := db.QueryRow("SELECT * FROM users WHERE name=$1", username)

	// get stored details
	storedCreds := &Credentials{}
	err := rows.Scan(&storedCreds.Id, &storedCreds.Username, &storedCreds.JoinDate, &storedCreds.Wins,
		&storedCreds.Losses, &storedCreds.Draws, &storedCreds.Hash)
	if err != nil {
		return nil, err
	}

	return storedCreds, nil
}
