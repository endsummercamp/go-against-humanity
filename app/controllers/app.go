package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/gorilla/websocket"
	"github.com/revel/revel"
)

func hashPassword(password string) string {
	hasher := sha256.New()
	io.WriteString(hasher, password)
	return hex.EncodeToString(hasher.Sum(nil))
}

var deck *models.Deck
var mm = &game.MatchManager{}
var ws = SocketServer{sync.Mutex{}, mm, map[int][]*websocket.Conn{}}
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

	if !mm.IsJoinable(id) {
		c.Flash.Error(fmt.Sprintf("Unable to join %d. The match doesn't exists, is already started or already ended.", id))
		c.FlashParams()
		return c.Redirect(App.Matches)
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
		c.Flash.Error("Cannot join an unjoined match")
		c.FlashParams()
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

	black_card := deck.NewRandomBlackCard()
	white_card := deck.NewRandomWhiteCard()

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

func (c App) PickCard() revel.Result {
	user := c.connected()

	if user == nil {
		return c.Redirect(App.Login)
	}

	matchId, err := strconv.Atoi(c.Params.Route.Get("matchId"))

	if err != nil {
		return c.NotFound("Invalid MatchId")
	}

	if !mm.UserJoined(matchId, user) {
		return c.Forbidden("Vbb.")
	}

	cardId, err := strconv.Atoi(c.Params.Route.Get("cardId"))

	if err != nil {
		return c.NotFound("Invalid CardId")
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
			card = c
			foundId = i
			break
		}
	}

	if foundId == -1 {
		return c.NotFound("Card not found.")
	}

	player.Cards = append(player.Cards[:foundId], player.Cards[foundId+1:]...)

	result := round.AddCard(card)

	if !result {
		return c.Forbidden("Already played")
	}

	return c.RenderJSON(nil)
}

func (c App) MatchNewBlackCard() revel.Result {
	user := c.connected()
	if user == nil {
		return c.Redirect(App.Login)
	}

	if !user.IsAdmin() {
		return c.Forbidden("Not allowed.")
	}

	matchId, err := strconv.Atoi(c.Params.Route.Get("matchId"))
	if err != nil {
		return c.NotFound("Invalid MatchId")
	}

	match := mm.GetMatchByID(matchId)

	if match == nil {
		return c.NotFound("Match not found")
	}

	card := match.NewBlackCard()
	if card == nil {
		/* ... */
	}

	msg := Event{
		Name:     "new_black",
		NewCard:  card,
		Duration: 5, // Timeout in seconds
	}

	round := match.GetRound()
	round.TimeFinishPick = time.Now()

	go func() {
		time.Sleep(time.Duration(msg.Duration) * time.Second)
		match.State = models.MATCH_VOTING
		msg := Event{
			Name: "voting",
		}
		ws.BroadcastToRoom(matchId, msg)

		for _, card := range round.GetChoices() {
			msg := Event{
				Name:    "new_white",
				NewCard: card,
			}
			ws.BroadcastToRoom(matchId, msg)
			time.Sleep(time.Second)
		}
	}()

	ws.BroadcastToRoom(matchId, msg)
	return c.RenderJSON(true)
}

func (c App) EndVoting() revel.Result {
	user := c.connected()

	if user == nil {
		return c.Redirect(App.Login)
	}

	matchId, err := strconv.Atoi(c.Params.Route.Get("matchId"))

	if err != nil {
		return c.NotFound("Invalid MatchId")
	}

	match := mm.GetMatchByID(matchId)
	if match == nil {
		return c.NotFound("Match not found.")
	}

	if match.State != models.MATCH_VOTING {
		return c.Forbidden("Unable to end voting because voting hasn't started yet.")
	}

	match.EndVote()
	msg := Event{
		Name: "show_results",
	}
	ws.BroadcastToRoom(matchId, msg)

	for _, player := range match.Players {
		if len(player.Cards) < 10 {
			player.Cards = append(player.Cards, match.Deck.NewRandomWhiteCard())
		}
	}

	return c.Render()
}

func (c App) VoteCard() revel.Result {
	user := c.connected()

	if user == nil {
		return c.Redirect(App.Login)
	}

	if user.UserType != models.JurorType {
		return c.Forbidden("Only Jurors can cast a vote!")
	}

	matchId, err := strconv.Atoi(c.Params.Route.Get("matchId"))
	if err != nil {
		return c.NotFound("Invalid matchId parameter")
	}

	cardId, err := strconv.Atoi(c.Params.Route.Get("cardId"))
	if err != nil {
		return c.NotFound("Invalid cardId parameter")
	}

	match := mm.GetMatchByID(matchId)

	if match == nil {
		return c.NotFound("Match not found")
	}

	round := match.GetRound()

	if round == nil {
		return c.NotFound("Round not found")
	}

	var card *models.WhiteCard = nil

	for _, c := range round.GetChoices() {
		if c.Id == cardId {
			card = c
			break
		}
	}

	if card == nil {
		return c.NotFound("Card not found")
	}

	if match.State != models.MATCH_VOTING {
		return c.Forbidden("Voting disallowed")
	}

	var juror *models.Juror = nil

	for _, j := range match.Jury {
		if j.User.Id == user.Id {
			juror = &j
			break
		}
	}

	if juror == nil {
		return c.NotFound("Juror not found! Are you a Juror in this match?")
	}

	log.Printf("User: %#v\n", user)
	log.Printf("Juror: %#v\n", juror.User)

	for _, j := range round.Voters {
		if j.User.Id == juror.User.Id {
			return c.Forbidden("Cannot vote twice.")
		}
	}

	round.Voters = append(round.Voters, *juror)


	// Cast vote
	round.Wcs[card] = append(round.Wcs[card], *juror)

	totals := []Total{}

	for card, jury := range round.Wcs {
		totals = append(totals, Total{
			ID:    card.Id,
			Votes: len(jury),
		})
	}

	ws.BroadcastToRoom(matchId, Event{
		Name:   "vote_cast",
		Totals: totals,
	})

	return c.Render()
}
