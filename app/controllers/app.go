package controllers

import (
	"encoding/json"
	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/revel/revel"
	"log"
	"os"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func hashPassword(password string) string {
	hasher := sha256.New()
	io.WriteString(hasher, password)
	return hex.EncodeToString(hasher.Sum(nil))
}

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
	c.ViewArgs["failed"] = c.Params.Get("failed") != ""
	c.ViewArgs["registered"] = c.Params.Get("registered") != ""
	return c.Render()
}

func (c App) PostLogin() revel.Result {
	username := c.Params.Form.Get("username")
	password := c.Params.Form.Get("password")
	user := User{}
	pwhash := hashPassword(password)
	err := DbMap.SelectOne(&user, "SELECT * FROM users WHERE username=? AND pwhash=?", username, pwhash)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Redirect("/login?failed=1")
		} else {
			panic(err)
		}
	}
	fmt.Printf("%#v\n", user)
	return c.Redirect("/")
}

func (c App) Register() revel.Result {
	return c.Render()
}

func (c App) PostRegister() revel.Result {
	username := c.Params.Form.Get("username")
	password := c.Params.Form.Get("password")
	count, err := DbMap.SelectInt("SELECT COUNT(*) FROM users WHERE username=?", username)
	if err != nil {
		panic(err)
	}
	if count != 0 {
		c.ViewArgs["error"] = "Another user with that username already exists."
		c.Render()
	}
	user := User{
		Username: username,
		PwHash:   hashPassword(password),
	}
	err = DbMap.Insert(&user)
	if err != nil {
		panic(err)
	}
	return c.Redirect("/login?registered=1")
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