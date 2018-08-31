package controllers

import (
	"encoding/json"
	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/revel/revel"
	"log"
	"os"
)

var deck *game.Deck

type App struct {
	*revel.Controller
	deck *game.Deck
}

func (c App) initDeck() {
	deck = new(game.Deck)
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Login() revel.Result {
	return c.Render()
}

func (c App) NewRound() revel.Result {
	generateDeck()

	return c.RenderJSON(deck)
}

func (c App) GetDeck() revel.Result {
	return c.RenderJSON(deck)
}

func (c App) Card() revel.Result {
	//generateDeck()

	c.ViewArgs["deck_name"] = deck.Name
	black_card := (*game.NewRandomCardFromDeck(game.BLACK_CARD, deck)).(game.BlackCard)
	white_card := (*game.NewRandomCardFromDeck(game.WHITE_CARD, deck)).(game.WhiteCard)

	c.ViewArgs["cards"] = []game.Card{white_card, black_card}

	return c.Render()
}

func deckAllowed(deckName string) bool {
	switch deckName {
	case "ita-original-sfoltita":
	case "ita-espansione":
	case "ita-HACK":
		return true
	}
	return false
}

func generateDeck() {
	if deck != nil {
		return
	}

	f, err := os.OpenFile("./cards/json-against-humanity/full.md.json", os.O_RDONLY, 755)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(f)
	var v game.DeckData
	if err := decoder.Decode(&v); err != nil {
		log.Fatal(err)
	}

	var whitecards []game.Card
	var blackcards []game.Card

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

	deck = &game.Deck{
		"test",
		blackcards,
		whitecards,
		nil,
	}
}