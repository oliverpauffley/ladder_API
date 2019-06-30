package main

import (
	"database/sql"
	"github.com/oliverpauffley/chess_ladder/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Mockdb struct {
	users map[string]models.CredentialsInternal
}

func (db Mockdb) CreateUser(username, password string) error {
	var mockcredentials models.CredentialsInternal
	mockcredentials.Username = username
	mockcredentials.Hash, _ = bcrypt.GenerateFromPassword([]byte(password), 8)
	db.users[username] = mockcredentials
	return nil
}

func (db Mockdb) QueryByName(username string) (models.CredentialsInternal, error) {
	mockcredentials, exists := db.users[username]
	if exists == false {
		return mockcredentials, sql.ErrNoRows
	}
	return mockcredentials, nil
}

func (db Mockdb) QueryById(id int) (models.CredentialsExternal, error) {
	for _, entry := range db.users {
		if entry.Id == id {
			print(id)
			user := db.users[entry.Username]
			userCredentials := models.CredentialsExternal{Id: user.Id, Username: user.Username, JoinDate: user.JoinDate.Round(time.Hour),
				Role: user.Role, Wins: user.Wins, Losses: user.Losses, Draws: user.Draws}
			return userCredentials, nil
		}
	}
	return models.CredentialsExternal{}, sql.ErrNoRows
}
