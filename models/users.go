package models

import (
	"github.com/pkg/errors"
	. "golang.org/x/crypto/bcrypt"
	"time"
)

// credentials struct to store user information
type CredentialsInternal struct {
	Id       int       `db:"id"`
	Username string    `db:"name"`
	Email    string    `db:"email"`
	JoinDate time.Time `db:"join_date"`
	Role     int       `db:"role"`
	Wins     int       `db:"wins"`
	Losses   int       `db:"losses"`
	Draws    int       `db:"draws"`
	Hash     []byte    `db:"hash"`
}

type CredentialsExternal struct {
	Id       int       `db:"id"`
	Username string    `db:"name"`
	Email    string    `db:"email"`
	JoinDate time.Time `db:"join_date"`
	Role     int       `db:"role"`
	Wins     int       `db:"wins"`
	Losses   int       `db:"losses"`
	Draws    int       `db:"draws"`
}

func (db *DB) CreateUser(username, email, password string) error {
	// salt and hash the password using bcrypt. salt set to 8
	hashedPassword, err := GenerateFromPassword([]byte(password), 8)
	if err != nil {
		// there is something wrong with the password hashing
		return err
	}

	// insert new user into db
	_, err = db.Query("INSERT INTO users (name, hash, email) VALUES ($1, $2, $3)",
		username, string(hashedPassword), email)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) QueryByEmail(email string) (CredentialsInternal, error) {
	sqlStatement := "SELECT id, name, email, join_date, wins, losses, draws, hash, role FROM users WHERE email=$1;"
	row := db.QueryRow(sqlStatement, email)
	// get stored details
	var storedCreds CredentialsInternal
	err := row.Scan(&storedCreds.Id, &storedCreds.Username, &storedCreds.Email, &storedCreds.JoinDate, &storedCreds.Wins,
		&storedCreds.Losses, &storedCreds.Draws, &storedCreds.Hash, &storedCreds.Role)
	if err != nil {
		return storedCreds, err
	}

	return storedCreds, nil
}

func (db *DB) QueryById(id int) (CredentialsExternal, error) {
	sqlStatement := "SELECT id, name, email, join_date, wins, losses, draws, role FROM users WHERE id=$1;"
	row := db.QueryRow(sqlStatement, id)
	// get stored details
	var storedCreds CredentialsExternal
	err := row.Scan(&storedCreds.Id, &storedCreds.Username, &storedCreds.Email, &storedCreds.JoinDate, &storedCreds.Wins,
		&storedCreds.Losses, &storedCreds.Draws, &storedCreds.Role)
	if err != nil {
		return storedCreds, err
	}

	return storedCreds, nil
}

func (db *DB) DeleteUser(id int) error {
	sqlStatement := "DELETE FROM users WHERE id = $1"
	row, err := db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	if num, _ := row.RowsAffected(); int(num) != 1 {
		return errors.New("Deleted more rows than was supposed to happen")
	}
	return nil
}
