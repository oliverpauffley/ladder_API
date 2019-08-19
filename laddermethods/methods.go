package laddermethods

import (
	"math"
)

type LadderMethod interface {
	AdjustRank(winnerPoints, loserPoints int, draw bool) (winnerNew, loserNew int)
	GetStartingValues() (points int)
}

type Elo struct {
	StartingPoints int
	ScaleFactor    float64
}

func (e Elo) GetStartingValues() (points int) {
	return e.StartingPoints
}

// function implements elo ranking system described at https://metinmediamath.wordpress.com/2013/11/27/how-to-calculate-the-elo-rating-including-example/
func (e Elo) AdjustRank(winnerPoints, loserPoints int, draw bool) (winnerNew, loserNew int) {
	// first transform both rankings
	winnerTransformed := math.Pow(10, float64(winnerPoints)/400)
	loserTransformed := math.Pow(10, float64(loserPoints)/400)

	// calculated expected score
	winnerExpected := winnerTransformed / (winnerTransformed + loserTransformed)
	loserExpected := loserTransformed / (winnerTransformed + loserTransformed)

	// update points based on result
	switch draw {
	case false:
		winnerNew = int(math.Round(float64(winnerPoints) + e.ScaleFactor*(1-winnerExpected)))
		loserNew = int(math.Round(float64(loserPoints) + e.ScaleFactor*(0-loserExpected)))
	case true:
		winnerNew = int(math.Round(float64(winnerPoints) + e.ScaleFactor*(0.5-winnerExpected)))
		loserNew = int(math.Round(float64(loserPoints) + e.ScaleFactor*(0.5-loserExpected)))
	}
	return winnerNew, loserNew
}
