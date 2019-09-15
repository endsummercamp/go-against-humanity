package utils

import (
	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/go-gorp/gorp"
	"github.com/labstack/echo"
)

type CustomContext struct {
	echo.Context
	Db *gorp.DbMap
	MatchManager *game.MatchManager
}