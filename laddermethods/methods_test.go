package laddermethods

import "testing"

func TestElo_AdjustRank(t *testing.T) {
	t.Run("updates user rank correctly on win", func(t *testing.T) {
		startingRanks := []int{2400, 2000}
		want := []int{2403, 1997}

		e := Elo{1000, 32}

		newWinner, newLoser := e.AdjustRank(startingRanks[0], startingRanks[1], false)

		if newWinner != want[0] || newLoser != want[1] {
			t.Errorf("Did not get the correct new ranks, got %d, %d, wanted %d, %d",
				newWinner, newLoser, want[0], want[1])
		}
	})
	t.Run("updates user rank correctly on a draw", func(t *testing.T) {
		startingRanks := []int{2400, 2000}
		want := []int{2387, 2013}

		e := Elo{1000, 32}

		newWinner, newLoser := e.AdjustRank(startingRanks[0], startingRanks[1], true)

		if newWinner != want[0] || newLoser != want[1] {
			t.Errorf("Did not get the correct new ranks, got %d, %d, wanted %d, %d",
				newWinner, newLoser, want[0], want[1])
		}
	})
}
