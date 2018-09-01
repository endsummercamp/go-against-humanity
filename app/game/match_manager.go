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

	new_match := models.NewMatch(mm.counter, []models.Player{})
	mm.matches = append(mm.matches, new_match)
	return new_match
}

func (mm *MatchManager) IsJoinable(id int) bool {
	return mm.GetMatchByID(id) != nil
}

func (mm *MatchManager) JoinMatch(id int, user *models.User) bool {
	match := mm.GetMatchByID(id)

	if match == nil {
		return false
	}

	player := models.Player{
		user,
		0,
		[]models.Card{},
	}

	match.Players = append(match.Players, player)

	return true
}

func (mm *MatchManager) GetMatchByID(id int) *models.Match {
	for _, m := range mm.matches {
		if m.Id == id {
			return m
		}
	}
	return nil
}