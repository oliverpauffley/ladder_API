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

type LadderUser struct {
	Id          int `database:"id"`
	UserId      int `database:"user_id"`
	LadderId    int `database:"ladder_id"`
	Rank        int `database:"rank"`
	HighestRank int `database:"highest_rank"`
	Points      int `database:"points"`
}

type LadderMethod interface {
	AdjustRank(Winner, Loser LadderUser) error
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

// get ladder from its HashId
func (db *DB) GetLadderFromHashId(hashId string) (Ladder, error) {
	// find ladder from its hash
	sqlStatement := "SELECT id, name, owner, hashid, method FROM ladders WHERE hashid = $1"
	row := db.QueryRow(sqlStatement, hashId)
	var ladder Ladder
	err := row.Scan(&ladder.Id, &ladder.Name, &ladder.Owner, &ladder.HashId, &ladder.Method)
	if err != nil {
		return Ladder{}, err
	}
	return ladder, nil
}

func (db *DB) JoinLadder(ladderId, userId int, method string) error {
	startingPoints := 0
	if method == "elo" {
		startingPoints = 1000
	}
	sqlStatement := "INSERT INTO ladders_users (ladder_id, user_id, points) VALUES ($1, $2, $3)"
	_, err := db.Exec(sqlStatement, ladderId, userId, startingPoints)
	if err != nil {
		return err
	}
	return nil
}

// TODO delete ladder
//  needs to delete all of the ladders_users references too

// TODO change ladder owner
