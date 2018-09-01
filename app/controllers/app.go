package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/gorilla/websocket"
	"github.com/revel/revel"
	"io"
	"log"
	"strconv"
)

func hashPassword(password string) string {
	hasher := sha256.New()
	io.WriteString(hasher, password)
	return hex.EncodeToString(hasher.Sum(nil))
}

var deck *models.Deck
var mm = &game.MatchManager{}
var ws = SocketServer{mm, map[int][]*websocket.Conn{}}
var _ = ws.Start()

type App struct {
	*revel.Controller
	deck *models.Deck
}

func (c App) connected() *models.User {
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

func (c App) isAdmin() bool {
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username).IsAdmin()
	}
	return false
}

func (c App) initDeck() {
	deck = new(models.Deck)
}

func (c App) getUser(username string) *models.User {
	user := models.User{}
	DbMap.SelectOne(&user, "SELECT * FROM users WHERE username=?", username)

	return &user
}

func (c App) Index() revel.Result {
	user := c.connected()
	if user == nil {
		return c.Redirect(App.Login)
	}

	c.ViewArgs["user"] = user

	return c.Render()
}

func (c App) Login() revel.Result {
	if c.connected() != nil {
		return c.Redirect(App.Index)
	}
	return c.Render()
}

func (c App) Matches() revel.Result {
	if c.connected() == nil {
		return c.Redirect(App.Login)
	}

	c.ViewArgs["matches"] = mm.GetMatches()
	return c.Render()
}

func (c App) PostLogin() revel.Result {
	username := c.Params.Form.Get("username")
	password := c.Params.Form.Get("password")
	user := models.User{}
	pwhash := hashPassword(password)
	err := DbMap.SelectOne(&user, "SELECT * FROM users WHERE username=? AND pwhash=?", username, pwhash)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.Flash.Error("Invalid username or password")
			c.FlashParams()
			return c.Redirect(App.Login)
		} else {
			panic(err)
		}
	}
	c.Session["user"] = string(user.Username)
	fmt.Printf("%#v\n", user)
	return c.Redirect(App.Login)
}

func (c App) Register() revel.Result {
	if c.connected() != nil {
		return c.Redirect(App.Index)
	}
	return c.Render()
}

func (c App) Logout() revel.Result {
	if c.connected() == nil {
		return c.Redirect(App.Login)
	}

	c.Flash.Success("Logged out successfully")
	c.Session = make(revel.Session)
	return c.Redirect(App.Login)
}

func (c App) JoinMatch(id int) revel.Result {
	user := c.connected()

	if user == nil {
		return c.Redirect(App.Login)
	}

	if !mm.IsJoinable(id){
		c.Flash.Error(fmt.Sprintf("Unable to join %d. The match doesn't exists or already ended.", id))
	}

	mm.JoinMatch(id, user)
	return c.Redirect(fmt.Sprintf("/match/%d", id))
}

func (c App) Match(id int) revel.Result {
	user := c.connected()

	if user == nil {
		return c.Redirect(App.Login)
	}

	if !mm.UserJoined(id, user) {
		return c.Redirect(App.Matches)
	}

	c.ViewArgs["user"] = user
	c.ViewArgs["match_id"] = id

	return c.Render()
}

func (c App) PostRegister() revel.Result {
	username := c.Params.Form.Get("username")
	password := c.Params.Form.Get("password")
	usertype := c.Params.Form.Get("user_type")

	user := models.User{
		Username: username,
		PwHash:   hashPassword(password),
	}

	if usertype == "player" {
		user.UserType = models.PlayerType
	} else {
		user.UserType = models.JurorType
	}

	count, err := DbMap.SelectInt("SELECT COUNT(*) FROM users WHERE username=?", username)
	if err != nil {
		log.Panic(err)
	}
	if count != 0 {
		c.Flash.Error("Another user with that username already exists.")
		c.FlashParams()
		return c.Redirect(App.Login)
	}
	err = DbMap.Insert(&user)
	if err != nil {
		panic(err)
	}
	c.Flash.Success("Registration completed! You may now login.")
	c.FlashParams()
	return c.Redirect(App.Login)
}

func (c App) NewMatch() revel.Result {
	if !c.isAdmin() {
		return c.Forbidden("Unauthorized.")
	}
	match := mm.NewMatch()

	c.Flash.Success(fmt.Sprintf("New Match created succesfully! (ID: %d)", match.Id))
	c.FlashParams()

	return c.Redirect(App.Admin)
}

func (c App) Admin() revel.Result {
	return c.Render()
}

func (c App) AdminUsers() revel.Result {
	user := models.User{}
	userlist, err := DbMap.Select(&user, "SELECT * FROM users")

	if err != nil {
		log.Fatal(err)
	}

	c.ViewArgs["users"] = userlist

	return c.Render()
}

func (c App) GetDeck() revel.Result {
	if !c.isAdmin() {
		return c.Forbidden("Unauthorized.")
	}

	return c.RenderJSON(deck)
}

func (c App) Card() revel.Result {
	if !c.isAdmin() {
		return c.Forbidden("Unauthorized.")
	}

	black_card := game.NewRandomCardFromDeck(models.BLACK_CARD, deck)
	white_card := game.NewRandomCardFromDeck(models.WHITE_CARD, deck)

	c.ViewArgs["cards"] = []models.Card{white_card, black_card}

	return c.Render()
}

func (c App) MyCards() revel.Result {
	user := c.connected()

	if user == nil {
		return c.Redirect(App.Login)
	}

	matchIdStr := c.Params.Query.Get("match_id")
	if matchIdStr == "" {
		return c.NotFound("The 'match_id' parameter is required.")
	}
	matchId, err := strconv.Atoi(matchIdStr)
	if err != nil {
		// Todo: implement
		panic(err)
	}
	if !mm.UserJoined(matchId, user) {
		return c.Redirect(App.Matches)
	}
	match := mm.GetMatchByID(matchId)
	matchPlayer := match.GetPlayerByID(user.Id)
	if matchPlayer == nil {
		c.Flash.Error("No such player")
		c.FlashParams()
		return c.Redirect(App.Matches)
	}
	return c.RenderJSON(matchPlayer.Cards)
}

func (c App) PickCard(matchId int, cardId int) revel.Result {
	user := c.connected()

	if user == nil {
		return c.Redirect(App.Login)
	}

	if !mm.UserJoined(matchId, user) {
		return c.Forbidden("Vbb.")
	}

	match := mm.GetMatchByID(matchId)
	if match == nil {
		return c.NotFound("Match not found.")
	}

	round := match.GetRound()

	if round == nil {
		return c.Forbidden("Can't play this card right now (no rounds available).")
	}

	if match.State != models.MATCH_PLAYBALE {
		return c.Forbidden("Can't play this card at this time.")
	}

	player := match.GetPlayerByID(user.Id)

	foundId := -1
	var card *models.WhiteCard = nil
	for i, c := range player.Cards {
		if c.Id == cardId {
			card = &c
			foundId = i
			break
		}
	}

	if foundId == -1 {
		return c.NotFound("Card not found.")
	}

	player.Cards = append(player.Cards[:foundId], player.Cards[foundId+1:]...)

	round.AddCard(card)

	return c.RenderJSON(nil)
}