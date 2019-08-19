package models

import "time"

type Game struct {
	Winner   int       `json:"winner"`
	Loser    int       `json:"loser"`
	Draw     bool      `json:"draw"`
	Date     time.Time `json:"date"`
	LadderId int       `json:"ladder_id"`
}

func (db *DB) AddGame(game Game) error {
	sqlStatement := "INSERT INTO games (winner, loser, draw, ladder_id) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(sqlStatement, game.Winner, game.Loser, game.Draw, game.LadderId)
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) GetResults(userId int) (wins, losses, draws int, err error) {
	// get wins draws and losses for this user
	sqlStatement := "SELECT" +
		"(SELECT COUNT(winner) AS wins FROM games WHERE winner=$1 AND draw=false)," +
		"(SELECT COUNT(loser) as losses FROM games WHERE loser=$1 AND draw=false)," +
		"(SELECT COUNT(winner) AS draws FROM games WHERE winner=$1 AND draw=true)"
	row := db.QueryRow(sqlStatement, userId)
	err = row.Scan(&wins, &losses, &draws)
	if err != nil {
		return 0, 0, 0, nil
	}
	return wins, losses, draws, nil
}
