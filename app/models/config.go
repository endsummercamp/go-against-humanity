package models

type Config struct {
	General ConfigGeneral
}

type ConfigGeneral struct {
	Decks []string
}
