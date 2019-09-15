package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func Logout(c echo.Context) error {
	s, _ := session.Get("session", c)
	s.Values = map[interface{}]interface{}{}
	_ = s.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
