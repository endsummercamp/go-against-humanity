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
	Ws           SocketServer
}

func (w *WebApp) GetUserByUsername(username string) *models.User {
	var ret *models.User

	err := w.Db.SelectOne(ret, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		return nil
	}
	return ret
}
