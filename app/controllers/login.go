package controllers

import (
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/gorilla/sessions"
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

	var flashData = data.FlashData{}

	if len(s.Flashes("success")) > 0 {
		flashData.Success = s.Flashes("success")[0].(string)
	}

	if len(s.Flashes("error")) > 0 {
		flashData.Error = s.Flashes("error")[0].(string)
	}

	return c.Render(http.StatusOK, "Login.html", data.LoginPageData{
		Flash: flashData,
	})
}

func DoLogin(c echo.Context) error {
	s, _ := session.Get("session", c)
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	username := c.FormValue("username")
	pwhash := utils.HashPassword(c.FormValue("password"))

	cc := c.(*utils.CustomContext)
	if cc.Db == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var user *models.User
	v, err := cc.Db.Select(&user, "SELECT * FROM users WHERE username=? AND pwhash=?", username, pwhash)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			s.AddFlash("Invalid username or password")
			return c.Render(http.StatusOK, "Login.html", data.LoginPageData{
				Flash: data.FlashData{Error: s.Flashes()[0].(string)},
			})
		} else {
			panic(err)
		}
	}

	s.Values["user"] = v[0].(*models.User);

	_ = s.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusSeeOther, "/")
}
