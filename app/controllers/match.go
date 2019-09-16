package controllers

import (
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)


func (w *WebApp) Matches(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	return c.Render(http.StatusOK, "Matches.html", data.MatchesPageData{
		Matches: w.MatchManager.GetMatches(),
		User: *w.GetUserByUsername(utils.GetUsername(c)),
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
	if !w.MatchManager.IsJoinable(matchId) {
		return c.Redirect(http.StatusFound, "/matches")
	}

	w.MatchManager.JoinMatch(matchId, user)
	match := w.MatchManager.GetMatchByID(matchId)

	return c.Render(http.StatusOK, "Match.html", data.MatchPageData{
		Match: *match,
		User: *w.GetUserByUsername(utils.GetUsername(c)),
	})
}