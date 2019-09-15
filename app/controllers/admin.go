package controllers

import (
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"net/http"

	"github.com/labstack/echo"
)

func Admin(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	cc := c.(*utils.CustomContext)
	user := cc.GetUserByUsername(utils.GetUsername(c))

	if user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if !user.Admin {
		return c.NoContent(http.StatusForbidden)
	}

	return c.Render(http.StatusOK, "Admin.html", data.AdminPageData{
		User: *user,
	})
}
