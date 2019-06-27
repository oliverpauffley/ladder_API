package models

import (
	. "golang.org/x/crypto/bcrypt"
	"time"
)

// credentials struct to store user information
type Credentials struct {
	Id       int       `db:"id"`
	Username string    `db:"name"`
	JoinDate time.Time `db:"join_date"`
	Role     int       `db:"role"`
	Wins     int       `db:"wins"`
	Losses   int       `db:"losses"`
	Draws    int       `db:"draws"`
	Hash     []byte    `db:"hash"`
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

func (db *DB) QueryByName(username string) (Credentials, error) {
	sqlStatement := "SELECT id, name, join_date, wins, losses, draws, hash, role FROM users WHERE name=$1;"
	row := db.QueryRow(sqlStatement, username)
	// get stored details
	var storedCreds Credentials
	err := row.Scan(&storedCreds.Id, &storedCreds.Username, &storedCreds.JoinDate, &storedCreds.Wins,
		&storedCreds.Losses, &storedCreds.Draws, &storedCreds.Hash, &storedCreds.Role)
	if err != nil {
		return storedCreds, err
	}

	return storedCreds, nil
}
