package controllers

import (
	"fmt"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func (w *WebApp) Matches(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	return c.Render(http.StatusOK, "Matches.html", data.MatchesPageData{
		Matches: w.MatchManager.GetMatches(),
		User:    *w.GetUserByUsername(utils.GetUsername(c)),
		Header: data.HeaderData{
			Title: "Matches",
			SubTitle: "Join a match from the following...",
		},
	})
}

func (w *WebApp) JoinLatestMatch(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	matches := w.MatchManager.GetMatches()
	if len(matches) == 0 {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	matchId := matches[len(matches) - 1].Id

	user := w.GetUserByUsername(utils.GetUsername(c))
	if !w.MatchManager.IsJoinable(matchId) {
		return c.Redirect(http.StatusFound, "/matches")
	}

	w.MatchManager.JoinMatch(matchId, user)
	match := w.MatchManager.GetMatchByID(matchId)

	return c.Render(http.StatusOK, "Match.html", data.MatchPageData{
		Match: *match,
		User:  *w.GetUserByUsername(utils.GetUsername(c)),
	})
}

func (w *WebApp) JoinMatch(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	matchId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	user := w.GetUserByUsername(utils.GetUsername(c))

	if !w.MatchManager.UserJoined(matchId, user) {
		if !w.MatchManager.IsJoinable(matchId) {
			return c.Redirect(http.StatusFound, "/matches")
		}

		joinResult := w.MatchManager.JoinMatch(matchId, user)

		if !joinResult {
			return c.Redirect(http.StatusTemporaryRedirect, "/matches");
		}
	}

	match := w.MatchManager.GetMatchByID(matchId)

	return c.Render(http.StatusOK, "Match.html", data.MatchPageData{
		Match: *match,
		User:  *w.GetUserByUsername(utils.GetUsername(c)),
	})
}

func (w *WebApp) MatchCards(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	matchId, err := strconv.Atoi(c.QueryParam("match_id"))
	user := w.GetUserByUsername(utils.GetUsername(c))

	// TODO: Check condition!
	if !w.MatchManager.IsJoinable(matchId) || !w.MatchManager.UserJoined(matchId, user) {
		return c.NoContent(http.StatusForbidden)
	}

	if err != nil {
		return err
	}

	match := w.MatchManager.GetMatchByID(matchId)
	matchPlayer := match.GetPlayerByID(user.Id)
	if matchPlayer == nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, matchPlayer.Cards)
}

func (w *WebApp) NewBlackCard(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))
	if user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if !user.Admin {
		return c.NoContent(http.StatusForbidden)
	}

	matchId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	match := w.MatchManager.GetMatchByID(matchId)
	if match == nil {
		return c.NoContent(http.StatusNotAcceptable)
	}

	if match.State != models.MATCH_SHOW_RESULTS &&
		match.State != models.MATCH_WAIT_USERS {
		return c.NoContent(http.StatusNotAcceptable);
	}

	card := match.NewBlackCard()
	if card == nil {
		/* ... */
	}

	duration := 20;
	expires :=  time.Now().Unix() + int64(duration);

	msg := Event{
		Name:     "new_black",
		NewCard:  card,
		Expires: expires,
		State: match.State,
	}

	round := match.GetRound()
	round.Expires = expires

	go func() {
		time.Sleep(time.Duration(expires - time.Now().Unix()) *time.Second)
		match.State = models.MATCH_VOTING

		// Removing cards from Player's deck
		/*for c, _ := range round.Wcs {
			for _, p := range match.Players {
				for _, uc :=
			}
		}*/

		w.Ws.BroadcastToRoom(matchId, Event{
			Name: "voting",
			State: match.State,
		})

		for _, card := range round.GetChoices() {
			time.Sleep(time.Second)
			w.Ws.BroadcastToRoom(matchId, Event{
				Name:    "new_white",
				NewCard: card,
				State: match.State,
			})
		}
	}()

	w.Ws.BroadcastToRoom(matchId, msg)
	return c.JSON(http.StatusOK, true)
}

func (w *WebApp) PickCard(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	matchId, err := strconv.Atoi(c.Param("match_id"))
	if err != nil {
		return err
	}
	cardId, err := strconv.Atoi(c.Param("card_id"))
	if err != nil {
		return err
	}
	user := w.GetUserByUsername(utils.GetUsername(c))


	// TODO: Check condition!
	if !w.MatchManager.IsJoinable(matchId) || !w.MatchManager.UserJoined(matchId, user) {
		return c.NoContent(http.StatusForbidden)
	}

	match := w.MatchManager.GetMatchByID(matchId)
	if match == nil {
		return c.String(http.StatusNotFound, "Match not found.")
	}

	round := match.GetRound()
	if round == nil {
		return c.String(http.StatusForbidden, "Can't play this card right now (no rounds available).")
	}

	if match.State != models.MATCH_PLAYBALE {
		return c.String(http.StatusForbidden, "Can't play this card at this time.")
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

	if card == nil {
		return c.String(http.StatusNotFound, "Card not found.")
	}

	for _, c := range player.Cards {
		log.Printf("P%d, C: %d\n", player.User.Id, c.Id)
	}

	player.Cards = append(player.Cards[:foundId], player.Cards[foundId+1:]...)

	log.Printf("-----------")

	for _, c := range player.Cards {
		log.Printf("P%d, C: %d\n", player.User.Id, c.Id)
	}

	result := round.AddCard(card)

	if !result {
		return c.String(http.StatusForbidden, "Already played")
	}

	return c.JSON(http.StatusOK, nil)
}

func (w *WebApp) VoteCard(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))

	if user.UserType != models.JurorType {
		return c.String(http.StatusForbidden, "Only Jurors can cast a vote!")
	}

	matchId, err := strconv.Atoi(c.Param("match_id"))
	if err != nil {
		return err
	}
	cardId, err := strconv.Atoi(c.Param("card_id"))
	if err != nil {
		return err
	}

	match := w.MatchManager.GetMatchByID(matchId)
	if match == nil {
		return c.String(http.StatusNotFound, "Match not found.")
	}

	round := match.GetRound()
	if round == nil {
		return c.String(http.StatusForbidden, "Can't play this card right now (no rounds available).")
	}

	var card *models.WhiteCard = nil

	for _, c := range round.GetChoices() {
		if c.Id == cardId {
			card = c
			break
		}
	}

	if card == nil {
		return c.String(http.StatusNotFound, "Card not found.")
	}

	// TODO!
	/*
	if match.State != models.MATCH_VOTING {
		return c.Forbidden("Voting disallowed")
	}
	*/

	var juror *models.Juror = nil

	for _, j := range match.Jury {
		if j.User.Id == user.Id {
			juror = j
			break
		}
	}

	if juror == nil {
		return c.String(http.StatusNotFound, "Juror not found! Are you a Juror in this match?")
	}

	log.Printf("User: %#v\n", user)
	log.Printf("Juror: %#v\n", juror.User)

	for _, j := range round.Voters {
		if j.User.Id == juror.User.Id {
			return c.String(http.StatusForbidden, "Cannot vote twice.")
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

	w.Ws.BroadcastToRoom(matchId, Event{
		Name:   "vote_cast",
		Totals: totals,
		State: match.State,
	})

	return c.JSON(http.StatusOK, nil)
}

func (w *WebApp) EndVoting(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))
	if user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	if !user.Admin {
		return c.NoContent(http.StatusForbidden)
	}


	matchId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.String(http.StatusNotFound, "Invalid MatchId")
	}

	match := w.MatchManager.GetMatchByID(matchId)
	if match == nil {
		return c.String(http.StatusNotFound, "Match not found.")
	}

	if match.State != models.MATCH_VOTING {
		return c.String(http.StatusForbidden, "Unable to end voting because voting hasn't started yet.")
	}

	match.EndVote()
	w.Ws.BroadcastToRoom(matchId, Event{
		Name: "show_results",
		State: match.State,
	})

	for _, player := range match.Players {
		// log.Printf("Range over players... (%d - %d)", player.User.Id, len(player.Cards))
		if len(player.Cards) < 10 {
			whitecard := match.Deck.NewRandomWhiteCard()
			// log.Printf("Whitecard player %d : %#v", &player.User.Id, whitecard)
			player.Cards = append(player.Cards, whitecard)
		}
	}

	// Choose winner
	round := match.GetRound()
	var totals []Total

	for card, jury := range round.Wcs {
		totals = append(totals, Total{
			ID: card.Id,
			Votes: len(jury),
		})
	}

	sort.Slice(totals, func(i, j int) bool {
		return totals[i].Votes < totals[j].Votes
	})

	var winner *models.Player
	var winningCard *models.WhiteCard

	if len(totals) > 0 {
		winningID := totals[0].ID
		for card := range round.Wcs {
			if card.Id != winningID {
				continue
			}
			winner = card.Owner
			winningCard = card
		}
	}

	if winner != nil && winningCard != nil {
		winner.User.Score++
		numUpdated, err := w.Db.Update(winner.User)
		if err != nil {
			log.Printf("Winner update: %s\n", err)
		}
		if numUpdated != 1 {
			log.Printf("Winner update: expected 1 update, got %d\n", numUpdated)
		}
		fmt.Printf("Winner: %s\n", winner.User.Username)
		w.Ws.BroadcastToRoom(matchId, Event{
			Name:   "winner",
			State: match.State,
			WinnerUsername: winner.User.Username,
			WinnerText: winningCard.Text,
		})
	}

	return c.JSON(http.StatusOK, nil)
}