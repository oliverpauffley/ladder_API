package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/oliverpauffley/chess_ladder/laddermethods"
)

// interface for all db methods for handlers to use
type Datastore interface {
	// User methods
	CreateUser(username, email, password string) error
	QueryByEmail(email string) (CredentialsInternal, error)
	QueryById(id int) (CredentialsExternal, error)
	DeleteUser(id int) error

	// Ladder methods
	AddLadder(name, method string, owner int) error
	GetLadder(ladderId int) (Ladder, error)
	GetLadderFromHashId(HashId string) (Ladder, error)
	JoinLadder(ladderId, userId int, method laddermethods.LadderMethod) error
	GetLadders(userId int) ([]LadderInfo, error)
	GetUserPoints(ladderId, userId int) (int, error)
	UpdatePoints(userId, ladderId, newPoints int) error

	// Game Methods
	AddGame(game Game) error
	GetResults(userId int) (wins, losses, draws int, err error)
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
