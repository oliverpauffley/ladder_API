package models

import (
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
