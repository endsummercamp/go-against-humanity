package main

import (
	"fmt"
	"github.com/ESCah/go-against-humanity/app/controllers"
	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
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
	_ = c.String(http.StatusInternalServerError, err.Error())
}

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.New("").Funcs(utils.FuncMap).ParseGlob("app/views/**/*.html")),
	}

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("SECRET"))))
	e.Renderer = t
	e.HTTPErrorHandler = customHTTPErrorHandler

	w := controllers.WebApp{
		MatchManager: new(game.MatchManager),
		Echo:         e,
		Db:           utils.InitDB(),
	}

	e.Static("/public", "public")
	e.GET("/login", w.Login)
	e.POST("/login", w.DoLogin)
	e.GET("/signup", w.SignUp)
	e.POST("/signup", w.DoSignUp)
	e.GET("/logout", w.Logout)

	e.GET("/admin", w.Admin)
	e.GET("/admin/users", w.AdminUsers)
	e.GET("/admin/matches/new", w.AdminNewMatch)
	e.GET("/matches", w.Matches)
	e.GET("/matches/join/:id", w.JoinMatch)

	e.GET("/", w.Index)
	e.Logger.Debug(e.Start(":1323"))
}
