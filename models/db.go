package models

import (
	"database/sql"
	_ "github.com/lib/pq"
)

// interface for all db methods for handlers to use
type Datastore interface {
	CreateUser(username, password string) error
	QueryByName(username string) (Credentials, error)
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
