package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/labstack/echo"

	"crypto/sha256"

	"encoding/hex"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

func hashPassword(password string) string {
	hasher := sha256.New()
	io.WriteString(hasher, password)
	return hex.EncodeToString(hasher.Sum(nil))
}

type Template struct {
	templates *template.Template
}

var funcMap = template.FuncMap{
	"replace": func(input, from, to string) string {
		return strings.Replace(input, from, to, -1)
	},
	"card_text": func(input models.Card) string {
		return input.GetText()
	},
	"card_dash": func(input models.Card) string {
		return strings.Replace(input.GetText(), "_", "<div class=\"long-dash\"></div>", -1)
	},
	"card_black": func(input models.Card) bool {
		return input.GetColor() == models.BLACK_CARD
	},
	"long_text": func(input string) bool {
		return len(input) > 100
	},
	"is_player": func(user models.User) bool {
		return user.UserType == models.PlayerType
	},
	"is_admin": func(user models.User) bool {
		return user.IsAdmin()
	},
	"format_date": func(date time.Time) string {
		return date.Format("2 Jan 2006, 15:04:01")
	},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
	c.String(http.StatusInternalServerError, err.Error())
}

type HeaderData struct {
	MoreStyles  []string
	MoreScripts []string
}

type FlashData struct {
	Success string
	Error   string
}

type LoginPageData struct {
	Header HeaderData
	Flash  FlashData
}

type SignupPageData struct {
	Header HeaderData
	Error  string
}

func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "Login.html", LoginPageData{})
}

func DoLogin(c echo.Context) error {
	s, _ := session.Get("session", c)
	s.Save(c.Request(), c.Response())
	return nil
}

func DoSignUp(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	user_type := c.FormValue("user_type")

	user := models.User{
		Username: username,
		PwHash:   hashPassword(password),
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
	return c.Render(http.StatusOK, "SignUp.html", SignupPageData{})
}

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("app/views/**/*.html")),
	}

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	e.Renderer = t
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Static("/public", "public")
	e.GET("/login", Login)
	e.POST("/login", DoLogin)
	e.GET("/signup", SignUp)
	e.POST("/signup", DoSignUp)
	e.Logger.Debug(e.Start(":1323"))
}
