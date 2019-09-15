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

	userVal := s.Values["user"]

	if userVal != nil {
		user := userVal.(models.User)
		if user.Username != "" {
			return c.Render(http.StatusOK, "Index.html", data.IndexPageData{User: user})
		}
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
