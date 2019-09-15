package utils

import (
	"github.com/go-gorp/gorp"
	"github.com/labstack/echo"
)

type CustomContext struct {
	echo.Context
	Db *gorp.DbMap
}