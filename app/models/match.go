package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	gc_log "github.com/denysvitali/gc_log"
)

type Match struct {
	Id        int
	Players   []*Player
	Jury      []*Juror
	CreatedOn time.Time
	Rounds    []Round
	State     MatchState
	Deck      *Deck
}

type MatchState int

const (
	MATCH_WAIT_USERS MatchState = iota
	MATCH_PLAYBALE
	MATCH_VOTING
	MATCH_SHOW_RESULTS
	MATCH_END
)

func NewMatch(id int, players []*Player) *Match {
	m := new(Match)
	m.Deck = nil
	m.Id = id
	m.Players = players
	m.CreatedOn = time.Now()
	m.State = MATCH_WAIT_USERS
	return m
}

var allowedDecks []string

func deckAllowed(deckName string) bool {
	for _, d := range allowedDecks {
		if d == deckName {
			return true
		}
	}
	return false
}

func (m *Match) NewDeck() {
	if m.Deck != nil {
		return
	}

	gc_log.Debug("Generating deck...")

	var conf Config
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		gc_log.Fatal(err)
	}

	gc_log.Debug(fmt.Sprintf("Decks allowed: %s", conf.General.Decks))

	allowedDecks = conf.General.Decks

	f, err := os.OpenFile("./cards/json-against-humanity/full.md.json", os.O_RDONLY, 755)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(f)
	var v DeckData
	if err := decoder.Decode(&v); err != nil {
		log.Fatal(err)
	}

	var whitecards []WhiteCard
	var blackcards []BlackCard

	i := 0

	for _, card := range v.White {
		if deckAllowed(card.Deck) {
			card.Id = i
			i++
			whitecards = append(whitecards, card)
		}
	}

	for _, card := range v.Black {
		if deckAllowed(card.Deck) && card.Pick == 1 {
			card.Id = i
			i++
			blackcards = append(blackcards, card)
		}
	}

	m.Deck = &Deck{
		blackcards,
		whitecards,
		nil,
	}
}

func (m *Match) GetPlayerByID(id int64) *Player {
	for _, player := range m.Players {
		if player.User.Id == id {
			return player
		}
	}
	return nil
}

func (m *Match) GetRound() *Round {
	if len(m.Rounds) == 0 {
		return nil
	}

	return &m.Rounds[len(m.Rounds)-1]
}

func (m *Match) NewBlackCard() *BlackCard {

	if m.State != MATCH_SHOW_RESULTS && m.State != MATCH_WAIT_USERS {
		return nil
	}

	m.State = MATCH_PLAYBALE

	blackCard := m.Deck.NewRandomBlackCard()
	if blackCard == nil {
		return nil
	}
	m.Rounds = append(m.Rounds, Round{
		BlackCard: blackCard,
		Wcs:       map[*WhiteCard][]Juror{},
		Mutex: 		sync.Mutex{},
		Voters: 	[]Juror{},
	})

	return blackCard
}

func (m *Match) EndVote() bool {
	if m.State != MATCH_VOTING {
		return false
	}

	m.State = MATCH_SHOW_RESULTS
	return true
}
func (m *Match) RemoveVote(round *Round, card *WhiteCard, juror *Juror) {
	found := -1
	for i, j := range round.Wcs[card] {
		if j.User.Id == juror.User.Id {
			found = i
			break
		}
	}

	if found == -1 {
		return
	}

	round.Wcs[card] = append(round.Wcs[card][:found], round.Wcs[card][found+1:]...)

}
