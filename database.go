package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// NewDatabase Create a new sql database, remember to db.Close() before end of program.
func NewDatabase() *sql.DB {
	db, err := sql.Open("mysql",
		"root:root@tcp(172.17.0.2:3306)/chess")
	if err != nil {
		panic(err)
	}
	return db
}
