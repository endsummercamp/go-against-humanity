package controllers

import (
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"log"
	"net/http"

	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/labstack/echo"
)

func DoSignUp(c echo.Context) error {
	s, err := session.Get("session", c)
	username := c.FormValue("username")
	password := c.FormValue("password")
	userType := c.FormValue("user_type")

	user := models.User{
		Username: username,
		PwHash:   utils.HashPassword(password),
	}

	if userType == "player" {
		user.UserType = models.PlayerType
	} else {
		user.UserType = models.JurorType
	}

	fmt.Printf("%#v\n", user)

	cc := c.(*utils.CustomContext)

	count, err := cc.Db.SelectInt("SELECT COUNT(*) FROM users WHERE username=?", username)
	if err != nil {
		log.Panic(err)
	}
	if count != 0 {
		return c.Render(http.StatusOK, "SignUp.html", data.SignupPageData{
			Flash: data.FlashData{
				Error: "Another user with that username already exists.",
			},
		})
	}
	err = cc.Db.Insert(&user)
	if err != nil {
		panic(err)
	}
	s.AddFlash("Registration completed! You may now login.", "success")

	/* c.String(http.StatusOK, fmt.Sprintf("U: %s, P: %s, T: %s", username, password,
	user_type)) */

	_ = c.Redirect(http.StatusSeeOther, "/")
	return nil
}

func SignUp(c echo.Context) error {
	return c.Render(http.StatusOK, "SignUp.html", data.SignupPageData{})
}
