package controllers

import (
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"net/http"

	"github.com/labstack/echo"
)

func (w *WebApp) Index(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))
	if user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(http.StatusOK, "Index.html", data.IndexPageData{
		User: *user,
	})
}
