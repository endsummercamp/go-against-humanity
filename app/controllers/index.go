package controllers

import (
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func (w *WebApp) Index(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))
	if user == nil {
		// This can occur if you reload the server and delete the db
		return w.Logout(c)
	}

	leaderboard, err := w.GetLeaderboard(user.Id)
	if err != nil {
		log.Printf("Failed to get leaderboard: %s\n", err)
	}
	log.Printf("Leaderboard: %#v\n", leaderboard)

	return c.Render(http.StatusOK, "Index.html", data.IndexPageData{
		User:        *user,
		Leaderboard: leaderboard,
		Header: data.HeaderData{
			Title:    "Home",
			SubTitle: "Welcome to Cards Against Humanity",
		},
	})
}
