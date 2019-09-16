package controllers

import (
	"fmt"
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/gorilla/sessions"
	"log"
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
		MaxAge:   3600,
		HttpOnly: true,
	}

	username := c.FormValue("username")

	cc := c.(*utils.CustomContext)
	if cc.Db == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	v, err := cc.Db.Select(models.User{}, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			s.AddFlash("Username not found")
			return c.Render(http.StatusOK, "Login.html", data.LoginPageData{
				Flash: data.FlashData{Error: s.Flashes()[0].(string)},
			})
		} else {
			panic(err)
		}
	}

	if len(v) == 0 || !utils.CheckPassword(v[0].(*models.User).PwHash, c.FormValue("password")) {
		s.AddFlash("Invalid username or password")
		return c.Render(http.StatusOK, "Login.html", data.LoginPageData{
			Flash: data.FlashData{Error: s.Flashes()[0].(string)},
		})
	}

	user := v[0].(*models.User)
	s.Values["user"] = username
	fmt.Printf("Logging in as %s\n", user.Username)

	err = s.Save(c.Request(), c.Response())
	if err != nil {
		log.Fatal(err)
	}
	return c.Redirect(http.StatusSeeOther, "/")
}
