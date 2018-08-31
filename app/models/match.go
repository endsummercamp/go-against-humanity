package models

import (
	"encoding/json"
	"log"
	"os"
	gc_log "github.com/denysvitali/gc_log"
)

type Match struct {
	Id	int
	Players []User
	deck 	*Deck
}

func NewMatch(id int, players []User) *Match {
	m := new(Match)
	m.deck = new(Deck)
	m.Id = id
	m.Players = players
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
	if m.deck != nil {
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

	m.deck = &Deck{
		blackcards,
		whitecards,
		nil,
	}
}