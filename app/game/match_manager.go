package game

import "github.com/ESCah/go-against-humanity/app/models"

type MatchManager struct {
	matches []*models.Match
	counter	int
}

func(mm *MatchManager) GetMatches() []*models.Match {
	return mm.matches
}

func(mm *MatchManager) NewMatch() *models.Match {
	mm.counter++

	new_match := models.NewMatch(mm.counter, []models.User{})
	mm.matches = append(mm.matches, new_match)
	return new_match
}