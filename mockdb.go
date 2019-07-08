package main

import (
	"database/sql"
	"github.com/oliverpauffley/chess_ladder/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Mockdb struct {
	users   map[string]models.CredentialsInternal
	ladders map[int]models.Ladder
}

func (db Mockdb) CreateUser(username, email, password string) error {
	var mockcredentials models.CredentialsInternal
	mockcredentials.Username = username
	mockcredentials.Hash, _ = bcrypt.GenerateFromPassword([]byte(password), 8)
	mockcredentials.Email = email
	db.users[username] = mockcredentials
	return nil
}

func (db Mockdb) QueryByEmail(email string) (models.CredentialsInternal, error) {
	for _, entry := range db.users {
		if entry.Email == email {
			user := db.users[entry.Username]
			userCredentials := models.CredentialsInternal{Id: user.Id, Username: user.Username, Email: user.Email, JoinDate: user.JoinDate.Round(time.Hour),
				Role: user.Role, Wins: user.Wins, Losses: user.Losses, Draws: user.Draws, Hash: user.Hash}
			return userCredentials, nil
		}
	}
	return models.CredentialsInternal{}, sql.ErrNoRows
}

func (db Mockdb) QueryById(id int) (models.CredentialsExternal, error) {
	for _, entry := range db.users {
		if entry.Id == id {
			print(id)
			user := db.users[entry.Username]
			userCredentials := models.CredentialsExternal{Id: user.Id, Username: user.Username, Email: user.Email, JoinDate: user.JoinDate.Round(time.Hour),
				Role: user.Role, Wins: user.Wins, Losses: user.Losses, Draws: user.Draws}
			return userCredentials, nil
		}
	}
	return models.CredentialsExternal{}, sql.ErrNoRows
}

func (db Mockdb) DeleteUser(id int) error {
	user, err := db.QueryById(id)
	if err != nil {
		return err
	}
	delete(db.users, user.Username)
	return nil
}

func (db Mockdb) AddLadder(name, method string, owner int) error {
	// find max key value
	key := 0
	for range db.ladders {
		if _, exists := db.ladders[key+1]; exists == true {
			key++
		}
	}
	// add new ladder
	newLadder := models.Ladder{Id: key, Name: name, Owner: owner, Method: method, HashId: "Ladder1"}
	db.ladders[key] = newLadder
	return nil
}
