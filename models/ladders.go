package models

import (
	"github.com/oliverpauffley/chess_ladder/laddermethods"
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
	Id       int `database:"id" json:"id"`
	UserId   int `database:"user_id" json:"user_id"`
	LadderId int `database:"ladder_id" json:"ladder_id"`
	Rank     int `database:"rank" json:"rank"`
	Points   int `database:"points" json:"points"`
}

type LadderRanks struct {
	Name   string `database:"name" json:"name"`
	UserId int    `database:"user_id" json:"user_id"`
	Rank   int    `database:"rank" json:"rank"`
	Points int    `database:"points" json:"points"`
}

type LadderInfo struct {
	LadderId int           `json:"ladder_id"`
	Name     string        `json:"name"`
	Owner    int           `json:"owner"`
	HashId   string        `json:"hash_id"`
	Players  []LadderRanks `json:"players"`
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

// add a user to a ladder
func (db *DB) JoinLadder(ladderId, userId int, method laddermethods.LadderMethod) error {
	sqlStatement := "INSERT INTO ladders_users (ladder_id, user_id, points) VALUES ($1, $2, $3)"
	_, err := db.Exec(sqlStatement, ladderId, userId, method.GetStartingValues())
	if err != nil {
		return err
	}
	return nil
}

// get all ladders that a user created
func (db *DB) GetLadders(userId int) ([]LadderInfo, error) {
	// Get all ladders that the user is in or where the user is the owner
	sqlStatement := "SELECT id, name, hashid, owner FROM ladders WHERE owner=$1 UNION SELECT ladders.id, ladders.name," +
		" ladders.hashid, ladders.owner FROM ladders_users join ladders ON ladders_users.ladder_id = ladders.id" +
		" WHERE ladders_users.user_id=$1"
	rows, err := db.Query(sqlStatement, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// start empty list for ladder information
	var userLadders []Ladder

	// scan over rows and add each ladder to the list of ladders
	// TODO Deal with nil return values
	for rows.Next() {
		var ladder Ladder
		if err := rows.Scan(&ladder.Id, &ladder.Name, &ladder.HashId, &ladder.Owner); err != nil {
			return nil, err
		}
		userLadders = append(userLadders, ladder)
	}

	// empty list to store ladder info with players attached
	var laddersWithPlayers []LadderInfo

	sqlStatement = "SELECT users.name, ladders_users.user_id, RANK () OVER (ORDER  BY ladders_users.points) rank," +
		" ladders_users.points FROM ladders_users JOIN users ON ladders_users.user_id = users.id WHERE ladders_users.ladder_id=$1"
	// Get players for each ladder
	for _, ladder := range userLadders {
		playerRow, err := db.Query(sqlStatement, ladder.Id)
		if err != nil {
			return nil, err
		}
		var playerList []LadderRanks
		// scan each player into the list
		for playerRow.Next() {
			var player LadderRanks
			err = playerRow.Scan(&player.Name, &player.UserId, &player.Rank, &player.Points)
			if err != nil {
				return nil, err
			}
			playerList = append(playerList, player)
		}

		ladderWithPlayers := LadderInfo{
			LadderId: ladder.Id,
			Name:     ladder.Name,
			Owner:    ladder.Owner,
			HashId:   ladder.HashId,
			Players:  playerList,
		}
		// add list to players to ladder
		laddersWithPlayers = append(laddersWithPlayers, ladderWithPlayers)
	}

	return laddersWithPlayers, nil
}

// TODO delete ladder
//  needs to delete all of the ladders_users references too

// TODO change ladder owner
