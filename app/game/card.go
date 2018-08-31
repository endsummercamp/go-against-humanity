package game

import (
	"log"
	"math/rand"
	"strings"
	"time"
)

type Card interface {
	GetText() string
	GetColor() CardColor
}

type BlackCard struct {
	Deck	string	`json:"deck"`
	Icon	string	`json:"icon"`
	Text	string	`json:"text"`
	Pick int	`json:"pick"`
}

type WhiteCard struct {
	Deck	string	`json:"deck"`
	Icon	string	`json:"icon"`
	Text	string	`json:"text"`
}

func (c WhiteCard) GetText() string {
	return c.Text
}

func (c WhiteCard) GetColor() CardColor {
	return WHITE_CARD
}

func (c BlackCard) GetText() string {
	return c.Text
}

func (c BlackCard) GetColor() CardColor {
	return BLACK_CARD
}

func NewCard(color CardColor, text string) Card {
	var c Card = nil
	switch color {
	case BLACK_CARD:
		c = &BlackCard{
			Deck: "hello",
			Icon: "default",
			Text: text,
			Pick: strings.Count(text, "_"),
		}
	case WHITE_CARD:
		c = &WhiteCard{
			Deck: "hello",
			Icon: "default",
			Text: text,
		}
	}
	return c
}

func NewRandomCardFromDeck(color CardColor, deck *Deck) *Card {
	rand.Seed(time.Now().Unix())

	var card *Card = nil

	switch color {
	case BLACK_CARD:
		if len(deck.Black_cards) == 0 {
			return nil
		}
		i := rand.Intn(len(deck.Black_cards))
		card = &deck.Black_cards[i]
		deck.Black_cards = append(deck.Black_cards[:i], deck.Black_cards[i+1:]...)
		return card
	case WHITE_CARD:
		if len(deck.White_cards) == 0 {
			return nil
		}
		i := rand.Intn(len(deck.White_cards))
		card = &deck.White_cards[i]
		log.Printf("%#v\n", card)
		deck.White_cards = append(deck.White_cards[:i], deck.White_cards[i+1:]...)
		return card
	default:
		return nil
	}
}