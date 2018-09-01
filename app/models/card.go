package models

import (
	"strings"
)

type Card interface {
	GetText() string
	GetColor() CardColor
	GetId()	int
}

type BlackCard struct {
	Id		int
	Deck	string	`json:"deck"`
	Icon	string	`json:"icon"`
	Text	string	`json:"text"`
	Pick int	`json:"pick"`
}

type WhiteCard struct {
	Id		int
	Deck	string	`json:"deck"`
	Icon	string	`json:"icon"`
	Text	string	`json:"text"`
	Owner	*Player
}

func (c WhiteCard) GetText() string {
	return c.Text
}

func (c WhiteCard) GetColor() CardColor {
	return WHITE_CARD
}

func (c WhiteCard) GetId() int {
	return c.Id
}

func (c BlackCard) GetText() string {
	return c.Text
}

func (c BlackCard) GetColor() CardColor {
	return BLACK_CARD
}

func (c BlackCard) GetId() int {
	return c.Id
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