package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

type DB struct {
	*sql.DB
}

// start new db with given postgres open string
func NewDB(dataSourceName string) (*DB, error) {

	// attempt to connect 10 times to db
	for attempt := 0; attempt > 10; attempt++ {
		db, err := sql.Open("postgres", dataSourceName)
		if err == nil {
			// db has connected, ping to confirm
			if err = db.Ping(); err == nil {
				return &DB{db}, nil
			}
		}
		time.Sleep(time.Second * 5)
	}
	// unable to connect after 10 tries, return err
	return &DB{}, errors.New("unable to connect to db")
}
