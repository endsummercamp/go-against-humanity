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

func (w *WebApp) playersList(match *models.Match) []models.Player {
	ret := make([]models.Player, len(match.Players))
	for i, player := range match.Players {
		ret[i] = *player
		// Redact sensitive data
		tmpUser := *(ret[i].User)
		tmpUser.PwHash = ""
		ret[i].User = &tmpUser
		ret[i].Cards = make([]*models.WhiteCard, 0)
	}
	return ret
}

func (w *WebApp) jurorsList(match *models.Match) []models.Juror {
	ret := make([]models.Juror, len(match.Jury))
	for i, juror := range match.Jury {
		ret[i] = *juror
		// Redact sensitive data
		ret[i].User.PwHash = ""
	}
	return ret
}

func (w *WebApp) Matches(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		log.Println("[Matches] Not logged in, redirecting to /login")
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	return c.Render(http.StatusOK, "Matches.html", data.MatchesPageData{
		Matches: w.MatchManager.GetMatches(),
		User:    *w.GetUserByUsername(utils.GetUsername(c)),
		Header: data.HeaderData{
			Title:    "Matches",
			SubTitle: "Join a match from the following...",
		},
	})
}

func (w *WebApp) JoinLatestMatch(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		log.Println("[JoinLatestMatch] Not logged in, redirecting to /login")
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	matches := w.MatchManager.GetMatches()
	if len(matches) == 0 {
		log.Println("[JoinLatestMatch] No active match, redirecting to /")
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	matchId := matches[len(matches)-1].Id

	user := w.GetUserByUsername(utils.GetUsername(c))
	if !w.MatchManager.IsJoinable(matchId) {
		log.Println("[JoinLatestMatch] Match is not joinable, redirecting to /matches")
		return c.Redirect(http.StatusFound, "/matches")
	}

	w.MatchManager.JoinMatch(matchId, user)
	match := w.MatchManager.GetMatchByID(matchId)

	w.Ws.BroadcastToRoom(matchId, Event{
		// The list of players has changed. Update it if you're watching it (i.e. are in projector view)
		Name:  "players_update",
		State: match.State,
		Leaderboard: w.playersList(match),
		Jury: w.jurorsList(match),
	})

	return c.Render(http.StatusOK, "Match.html", data.MatchPageData{
		Match: *match,
		User:  *w.GetUserByUsername(utils.GetUsername(c)),
	})
}

func (w *WebApp) JoinMatch(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		log.Println("[JoinMatch] Not logged in, redirecting to /login")
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	matchId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	user := w.GetUserByUsername(utils.GetUsername(c))

	if !w.MatchManager.UserJoined(matchId, user) {
		if !w.MatchManager.IsJoinable(matchId) {
			log.Println("[JoinMatch] Not joined and not joinable, redirecting to /matches")
			return c.Redirect(http.StatusFound, "/matches")
		}

		joinResult := w.MatchManager.JoinMatch(matchId, user)

		if !joinResult {
			log.Println("[JoinMatch] mm.JoinMatch failed, redirecting to /matches")
			return c.Redirect(http.StatusTemporaryRedirect, "/matches");
		}
	}

	match := w.MatchManager.GetMatchByID(matchId)

	w.Ws.BroadcastToRoom(matchId, Event{
		// The list of players has changed. Update it if you're watching it (i.e. are in projector view)
		Name:  "players_update",
		State: match.State,
		Leaderboard: w.playersList(match),
		Jury: w.jurorsList(match),
	})

	return c.Render(http.StatusOK, "Match.html", data.MatchPageData{
		Match: *match,
		User:  *w.GetUserByUsername(utils.GetUsername(c)),
	})
}

func (w *WebApp) MatchCards(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		log.Println("[MatchCards] Not logged in, redirecting to /login")
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	matchId, err := strconv.Atoi(c.QueryParam("match_id"))
	user := w.GetUserByUsername(utils.GetUsername(c))

	// TODO: Check condition!
	if !w.MatchManager.IsJoinable(matchId) || !w.MatchManager.UserJoined(matchId, user) {
		log.Println("[MatchCards] Either not joinable or not joined, returning 403")
		return c.NoContent(http.StatusForbidden)
	}

	if err != nil {
		return err
	}

	match := w.MatchManager.GetMatchByID(matchId)
	matchPlayer := match.GetPlayerByID(user.Id)
	if matchPlayer == nil {
		log.Println("[MatchCards] No such player, returning 500")
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, matchPlayer.Cards)
}

func (w *WebApp) NewBlackCard(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		log.Println("[NewBlackCard] Not logged in, redirecting to /login")
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))
	if user == nil {
		log.Println("[NewBlackCard] No such user, returning 500")
		return c.NoContent(http.StatusInternalServerError)
	}

	if !user.Admin {
		log.Println("[MatchCards] Not a user, returning 403")
		return c.NoContent(http.StatusForbidden)
	}

	matchId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("[MatchCards] Invalid id")
		return err
	}

	match := w.MatchManager.GetMatchByID(matchId)
	if match == nil {
		log.Println("[MatchCards] No such match, returning NotAcceptable")
		return c.NoContent(http.StatusNotAcceptable)
	}

	if match.State != models.MATCH_SHOW_RESULTS &&
		match.State != models.MATCH_WAIT_USERS {
		log.Println("[MatchCards] State doesn't allow dealing a new card")
		return c.NoContent(http.StatusNotAcceptable);
	}

	card := match.NewBlackCard()
	if card == nil {
		/* ... */
	}

	msg := Event{
		Name:    "new_black",
		NewCard: card,
		State:   match.State,
	}

	w.Ws.BroadcastToRoom(matchId, msg)
	return c.JSON(http.StatusOK, true)
}

func (w *WebApp) PickCard(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		log.Println("[PickCard] Not logged in, redirecting to /login")
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	matchId, err := strconv.Atoi(c.Param("match_id"))
	if err != nil {
		log.Println("[PickCard] Invalid param match_id")
		return err
	}
	cardId, err := strconv.Atoi(c.Param("card_id"))
	if err != nil {
		log.Println("[PickCard] Invalid param card_id")
		return err
	}
	user := w.GetUserByUsername(utils.GetUsername(c))

	// TODO: Check condition!
	if !w.MatchManager.IsJoinable(matchId) || !w.MatchManager.UserJoined(matchId, user) {
		log.Println("[PickCard] Neither joinable nor joined")
		return c.NoContent(http.StatusForbidden)
	}

	match := w.MatchManager.GetMatchByID(matchId)
	if match == nil {
		log.Println("[PickCard] No such match")
		return c.String(http.StatusNotFound, "Match not found.")
	}

	round := match.GetRound()
	if round == nil {
		log.Println("[PickCard] No rounds available")
		return c.String(http.StatusForbidden, "Can't play this card right now (no rounds available).")
	}

	if match.State != models.MATCH_PLAYBALE {
		log.Println("[PickCard] State doesn't allow picking a card")
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
		log.Println("[PickCard] Card not found")
		return c.String(http.StatusNotFound, "Card not found.")
	}

	/*
		for _, c := range player.Cards {
			log.Printf("P%d, C: %d\n", player.User.Id, c.Id)
		}
	*/

	player.Cards = append(player.Cards[:foundId], player.Cards[foundId+1:]...)

	/*
		log.Printf("-----------")

		for _, c := range player.Cards {
			log.Printf("P%d, C: %d\n", player.User.Id, c.Id)
		}
	*/

	result := round.AddCard(card)
	if !result {
		log.Println("[PickCard] Card already played")
		return c.String(http.StatusForbidden, "Already played")
	}

	// Show a white card with the player's name if it is not yet done.
	// Otherwise, go straight to voting.
	if len(round.Wcs) != len(match.Players) {
		w.Ws.BroadcastToRoom(matchId, Event{
			// The list of players has changed. Update it if you're watching it (i.e. are in projector view)
			Name:  "hidden_white_card",
			State: match.State,
			Username: player.User.Username,
		})
	} else {
		match.State = models.MATCH_VOTING

		// Removing cards from Player's deck
		/*for c, _ := range round.Wcs {
			for _, p := range match.Players {
				for _, uc :=
			}
		}*/

		w.Ws.BroadcastToRoom(matchId, Event{
			Name:  "voting",
			State: match.State,
		})

		for _, card := range round.GetChoices() {
			time.Sleep(time.Second)
			w.Ws.BroadcastToRoom(matchId, Event{
				Name:    "new_white",
				NewCard: card,
				State:   match.State,
			})
		}
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
		State:  match.State,
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
		Name:  "show_results",
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
			ID:    card.Id,
			Votes: len(jury),
		})
	}

	sort.Slice(totals, func(i, j int) bool {
		return totals[i].Votes > totals[j].Votes
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
			Name:           "winner",
			State:          match.State,
			WinnerUsername: winner.User.Username,
			WinnerText:     winningCard.Text,
		})
	}

	return c.JSON(http.StatusOK, nil)
}
