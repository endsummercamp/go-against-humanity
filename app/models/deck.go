package models

import (
	"log"
	"math/rand"
)

type Deck struct {
	Black_cards             	[]*BlackCard
	White_cards             	[]*WhiteCard
	LastExtractedCard 			*Card
}

type DeckMetadata struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Official    bool   `json:"official"`
}

type DeckData struct {
	Black []BlackCard	`json:"black"`
	White []WhiteCard	`json:"white"`
	Metadata map[string]DeckMetadata	 `json:"metadata"`
}

func RemoveBlackCard(cards []*BlackCard, index int) []*BlackCard {
	log.Printf("Removing %s...\n", cards[index].GetText())
	cards[len(cards)-1], (cards)[index] = (cards)[index], cards[len(cards)-1]
	cards = (cards)[:len(cards)-1]
	return cards
}


func RemoveWhiteCard(cards []*WhiteCard, index int) []*WhiteCard {
	log.Printf("Removing %s...\n", cards[index].GetText())
	cards[len(cards)-1], (cards)[index] = (cards)[index], cards[len(cards)-1]
	cards = (cards)[:len(cards)-1]
	return cards
}

func (deck *Deck) NewRandomWhiteCard() *WhiteCard {
	if len(deck.White_cards) == 0 {
		return nil
	}
	i := rand.Intn(len(deck.White_cards))
	card := &deck.White_cards[i]
	deck.White_cards = RemoveWhiteCard(deck.White_cards, i)
	return *card
}

func (deck *Deck) NewRandomBlackCard() *BlackCard {
	if len(deck.Black_cards) == 0 {
		return nil
	}
	i := rand.Intn(len(deck.Black_cards))
	card := &deck.Black_cards[i]
	log.Printf("Removing card %s\n", card)
	deck.Black_cards = RemoveBlackCard(deck.Black_cards, i)
	return *card
}