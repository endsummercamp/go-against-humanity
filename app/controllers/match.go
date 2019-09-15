package controllers

import (
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/labstack/echo"
	"net/http"
)

func Matches(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	cc := c.(*utils.CustomContext)
	return c.Render(http.StatusOK, "Matches.html", data.MatchesPageData{
		Matches: cc.MatchManager.GetMatches(),
		User: *cc.GetUserByUsername(utils.GetUsername(c)),
	})
}