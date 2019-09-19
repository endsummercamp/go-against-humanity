package models

type Card interface {
	GetText() string
	GetColor() CardColor
	GetId() int
}

type BlackCard struct {
	Id   int
	Deck string `json:"deck"`
	Icon string `json:"icon"`
	Text string `json:"text"`
	Pick int    `json:"pick"`
}

type WhiteCard struct {
	Id    int
	Deck  string  `json:"deck"`
	Icon  string  `json:"icon"`
	Text  string  `json:"text"`
	Owner *Player `json:"-"`
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
