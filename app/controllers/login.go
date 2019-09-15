package controllers

import (
	"net/http"

	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func Login(c echo.Context) error {
	s, _ := session.Get("session", c)
	if s.Values["user"] != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return c.Render(http.StatusOK, "Login.html", data.LoginPageData{})
}

func DoLogin(c echo.Context) error {
	s, _ := session.Get("session", c)
	/*s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}*/
	s.Values["user"] = models.User{
		Username: "test",
		Id:       1,
		Admin:    true,
	}
	s.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}
