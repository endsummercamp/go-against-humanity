package controllers

import (
	"fmt"
	"net/http"

	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/ESCah/go-against-humanity/app/models/data"
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/labstack/echo"
)

func DoSignUp(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	user_type := c.FormValue("user_type")

	user := models.User{
		Username: username,
		PwHash:   utils.HashPassword(password),
	}

	if user_type == "player" {
		user.UserType = models.PlayerType
	} else {
		user.UserType = models.JurorType
	}

	fmt.Printf("%#v\n", user)

	/*count, err := DbMap.SelectInt("SELECT COUNT(*) FROM users WHERE username=?", username)
	if err != nil {
		log.Panic(err)
	}
	if count != 0 {
		c.Flash.Error("Another user with that username already exists.")
		c.FlashParams()
		return c.Redirect(App.Login)
	}
	err = DbMap.Insert(&user)
	if err != nil {
		panic(err)
	}
	c.Flash.Success("Registration completed! You may now login.")
	c.FlashParams()

	/* c.String(http.StatusOK, fmt.Sprintf("U: %s, P: %s, T: %s", username, password,
	user_type)) */

	c.Redirect(http.StatusTemporaryRedirect, "/")
	return nil
}

func SignUp(c echo.Context) error {
	return c.Render(http.StatusOK, "SignUp.html", data.SignupPageData{})
}
