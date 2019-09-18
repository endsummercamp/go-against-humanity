package controllers

import (
	"fmt"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"net/http"

	"github.com/labstack/echo"
)

func (w *WebApp) Admin(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))

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

func (w *WebApp) AdminUsers(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))
	if user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if !user.Admin {
		return c.NoContent(http.StatusForbidden)
	}

	v, err := w.Db.Select(models.User{}, "SELECT * FROM users")
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var users []models.User
	for _, e := range v {
		users = append(users, *e.(*models.User))
	}

	return c.Render(http.StatusOK, "AdminUsers.html", data.AdminUsersPageData{
		User: *user,
		Users: users,
		Header: data.HeaderData{
			Title: "Users",
			// SubTitle: "List of the users",
		},
	})
}

func (w *WebApp) AdminNewMatch(c echo.Context) error {
	if !utils.IsLoggedIn(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user := w.GetUserByUsername(utils.GetUsername(c))
	if user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if !user.Admin {
		return c.NoContent(http.StatusForbidden)
	}

	m := w.MatchManager.NewMatch()
	return c.Redirect(http.StatusFound, fmt.Sprintf("/matches/join/%d", m.Id))
}