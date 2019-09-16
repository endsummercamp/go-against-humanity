package game

import (
	"github.com/ESCah/go-against-humanity/app/models"
)

type MatchManager struct {
	matches []*models.Match
	counter	int
}

func(mm *MatchManager) GetMatches() []*models.Match {
	return mm.matches
}

func(mm *MatchManager) NewMatch() *models.Match {
	mm.counter++

	newMatch := models.NewMatch(mm.counter, []*models.Player{})
	newMatch.NewDeck()
	mm.matches = append(mm.matches, newMatch)
	return newMatch
}

func (mm *MatchManager) IsJoinable(id int) bool {
	return mm.GetMatchByID(id) != nil
}

func (mm *MatchManager) JoinMatch(id int, user *models.User) bool {
	match := mm.GetMatchByID(id)

	if match == nil {
		return false
	}

	if match.State != models.MATCH_WAIT_USERS {
		return false
	}

	if user.UserType == models.PlayerType {
		for _, p := range match.Players {
			if p.User.Id == user.Id {
				return true
			}
		}

		var cards []*models.WhiteCard

		player := models.Player{}

		for i := 0; i < 10; i++ {
			card := match.Deck.NewRandomWhiteCard()
			card.Owner = &player
			cards = append(cards, card)
		}

		player.User = user
		player.Cards = cards

		match.Players = append(match.Players, &player)
	} else {
		for _, j := range match.Jury {
			if j.User.Id == user.Id {
				return true
			}
		}

		juror := models.Juror {
			User: user,
		}

		match.Jury = append(match.Jury, &juror)
	}

	return true
}

func (mm *MatchManager) UserJoined (id int, user *models.User) bool {
	match := mm.GetMatchByID(id)
	if match == nil {
		return false
	}

	for _, p := range match.Players {
		if p.User.Id == user.Id {
			return true
		}
	}

	for _, p := range match.Jury {
		if p.User.Id == user.Id {
			return true
		}
	}

	return false
}

func (mm *MatchManager) GetMatchByID(id int) *models.Match {
	for _, m := range mm.matches {
		if m.Id == id {
			return m
		}
	}
	return nil
}