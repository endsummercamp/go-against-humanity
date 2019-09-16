package controllers

import (
	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/go-gorp/gorp"
	"github.com/labstack/echo"
)

type WebApp struct {
	MatchManager *game.MatchManager
	Echo         *echo.Echo
	Db           *gorp.DbMap
}

func (w *WebApp) GetUserByUsername(username string) *models.User {
	res, err := w.Db.Select(models.User{}, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		return nil
	}
	if res != nil && len(res) == 1 {
		return res[0].(*models.User)
	}

	return nil
}
