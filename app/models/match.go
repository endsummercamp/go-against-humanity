package models

import (
	"encoding/json"
	"log"
	"os"
	gc_log "github.com/denysvitali/gc_log"
	"time"
)

type Match struct {
	Id           int
	Players      []Player
	Jury         []Juror
	CreatedOn    time.Time
	Rounds       []Round
	State		 MatchState
	Deck         *Deck
}

type MatchState int

const (
	MATCH_WAIT_USERS MatchState = iota
	MATCH_PLAYBALE
	MATCH_VOTING
	MATCH_SHOW_RESULTS
	MATCH_END
)

func NewMatch(id int, players []Player) *Match {
	m := new(Match)
	m.Deck = nil
	m.Id = id
	m.Players = players
	m.CreatedOn = time.Now()
	m.State = MATCH_WAIT_USERS
	return m
}

func deckAllowed(deckName string) bool {
	// TODO: Allow customization
	switch deckName {
	case "ita-original-sfoltita":
	case "ita-espansione":
	case "ita-HACK":
		return true
	}
	return false
}

func(m *Match) NewDeck(){
	if m.Deck != nil {
		return
	}

	gc_log.Debug("Generating deck...")

	f, err := os.OpenFile("./cards/json-against-humanity/full.md.json", os.O_RDONLY, 755)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(f)
	var v DeckData
	if err := decoder.Decode(&v); err != nil {
		log.Fatal(err)
	}

	var whitecards []Card
	var blackcards []Card

	for _, card := range v.White {
		if deckAllowed(card.Deck) {
			whitecards = append(whitecards, card)
		}
	}

	for _, card := range v.Black {
		if deckAllowed(card.Deck) {
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
			return &player
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