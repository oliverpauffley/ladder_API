package main

import (
	"database/sql"
	"github.com/oliverpauffley/chess_ladder/models"
	"golang.org/x/crypto/bcrypt"
)

type Mockdb struct {
	users map[string]models.Credentials
}

func (db Mockdb) CreateUser(username, password string) error {
	var mockcredentials models.Credentials
	mockcredentials.Username = username
	mockcredentials.Hash, _ = bcrypt.GenerateFromPassword([]byte(password), 8)
	db.users[username] = mockcredentials
	return nil
}

func (db Mockdb) QueryByName(username string) (models.Credentials, error) {
	mockcredentials, exists := db.users[username]
	if exists == false {
		return mockcredentials, sql.ErrNoRows
	}
	return mockcredentials, nil
}
