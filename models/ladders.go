package models

import (
	"github.com/speps/go-hashids"
)

type Ladder struct {
	Id     int    `database:"id"`
	Name   string `database:"name"`
	Owner  int    `database:"owner"`
	Method string `database:"method"`
	HashId string `database:"hashid"`
}

type LadderMethod interface {
	AdjustRank(Winner, Loser CredentialsInternal) error
}

// add new ladder
func (db *DB) AddLadder(name, method string, owner int) error {
	var id int
	sqlStatement := "INSERT INTO ladders (name, method, owner) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(sqlStatement, name, method, owner).Scan(&id)
	if err != nil {
		return err
	}

	//generate hashid from ladder id, uses github.com/speps/go-hashids
	hd := hashids.NewData()
	hd.Salt = "Secret Salt"
	hd.MinLength = 20
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return err
	}
	hashid, err := h.Encode([]int{int(id)})
	if err != nil {
		return err
	}

	// insert hashid into db
	sqlStatement = "UPDATE ladders SET hashid=$1 WHERE id=$2"
	_, err = db.Exec(sqlStatement, hashid, id)
	if err != nil {
		return err
	}
	return nil
}

// TODO delete ladder

// TODO change ladder owner
