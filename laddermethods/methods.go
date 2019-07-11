package laddermethods

type LadderMethod interface {
	AdjustRank(winnerPoints, loserPoints int) error
	GetStartingValues() (points int)
}

type Elo struct {
	StartingPoints int
}

func (e Elo) GetStartingValues() (points int) {
	return e.StartingPoints
}

// dummy function at the moment
func (e Elo) AdjustRank(winnerPoints, loserPoints int) error {
	return nil
}
