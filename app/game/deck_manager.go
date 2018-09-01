package game

import (
	"github.com/ESCah/go-against-humanity/app/models"
	"log"
	"math/rand"
	"fmt"
)

func RemoveCard(cards []models.Card, index int) []models.Card {
	log.Printf("Removing %s...\n", cards[index].GetText())
	cards[len(cards)-1], (cards)[index] = (cards)[index], cards[len(cards)-1]
	cards = (cards)[:len(cards)-1]
	return cards
}

func NewRandomCardFromDeck(color models.CardColor, deck *models.Deck) models.Card {
	var card models.Card

	fmt.Printf("%#v\n", *deck)
	switch color {
	case models.BLACK_CARD:
		if len(deck.Black_cards) == 0 {
			return nil
		}
		i := rand.Intn(len(deck.Black_cards))
		card = deck.Black_cards[i]
		log.Printf("Removing card %s\n", card)
		deck.Black_cards = RemoveCard(deck.Black_cards, i)
		return card
	case models.WHITE_CARD:
		if len(deck.White_cards) == 0 {
			return nil
		}
		i := rand.Intn(len(deck.White_cards))
		card = deck.White_cards[i]
		deck.White_cards = RemoveCard(deck.White_cards, i)
		return card
	default:
		return nil
	}
}