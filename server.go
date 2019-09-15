package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ESCah/go-against-humanity/app/controllers"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

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

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("app/views/**/*.html")),
	}

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	e.Renderer = t
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Static("/public", "public")
	e.GET("/login", controllers.Login)
	e.POST("/login", controllers.DoLogin)
	e.GET("/signup", controllers.SignUp)
	e.POST("/signup", controllers.DoSignUp)
	e.Logger.Debug(e.Start(":1323"))
}
