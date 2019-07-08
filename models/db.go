package models

import (
	"database/sql"
	_ "github.com/lib/pq"
)

// interface for all db methods for handlers to use
type Datastore interface {
	// User methods
	CreateUser(username, email, password string) error
	QueryByEmail(email string) (CredentialsInternal, error)
	QueryById(id int) (CredentialsExternal, error)
	DeleteUser(id int) error

	// Ladder methods
	AddLadder(name, method string, owner int) error
}

type DB struct {
	*sql.DB
}

// start new db with given postgres open string
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
