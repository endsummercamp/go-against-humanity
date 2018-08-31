package controllers

import (
	"encoding/json"
	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/revel/revel"
	"log"
	"os"
)

type App struct {
	*revel.Controller
	deck *game.Deck
}

func (c App) initDeck() {
	c.deck = new(game.Deck)
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Login() revel.Result {
	return c.Render()
}

func (c App) NewRound() revel.Result {
	f, err := os.OpenFile("./cards/json-against-humanity/full.md.json", os.O_RDONLY, 755)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(f)
	var v game.DeckData
	if err := decoder.Decode(&v); err != nil {
		log.Fatal(err)
	}

	whitecards := []game.Card{}
	blackcards := []game.Card{}

	for _, card := range v.White {
		whitecards = append(whitecards, card)
	}

	for _, card := range v.Black {
		blackcards = append(blackcards, card.WhiteCard)
	}

	c.deck = &game.Deck{
		"test",
		blackcards,
		whitecards,
		nil,
	}

	return c.RenderJSON(v)
}

func (c App) Card() revel.Result {

	c.ViewArgs["cards"] = []game.Card{game.NewCard(game.BLACK_CARD,
		"aaa")}

	c.ViewArgs["deck_name"] = c.deck.Name

	//game.NewRandomCardFromDeck(game.BLACK_CARD, c.deck)

	return c.Render()
}