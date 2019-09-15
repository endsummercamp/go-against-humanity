package controllers

import (
	"net/http"

	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func Index(c echo.Context) error {
	s, _ := session.Get("session", c)

	user_val := s.Values["user"]

	if user_val != nil {
		user := user_val.(models.User)
		if user.Username != "" {
			return c.Render(http.StatusOK, "Index.html", data.IndexPageData{User: user})
		}
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
