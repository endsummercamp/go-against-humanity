package controllers

import (
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func Index(c echo.Context) error {
	s, _ := session.Get("session", c)

	cc := c.(*utils.CustomContext)

	var user *models.User
	var username string
	if x, ok := s.Values["user"]; ok {
		username, _ = x.(string)
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user = cc.GetUserByUsername(username)

	if user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}


	return c.Render(http.StatusOK, "Index.html", data.IndexPageData{
		User: *user,
	})
}
