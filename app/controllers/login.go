package controllers

import (
	"net/http"

	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "Login.html", data.LoginPageData{})
}

func DoLogin(c echo.Context) error {
	s, _ := session.Get("session", c)
	s.Save(c.Request(), c.Response())
	return nil
}
