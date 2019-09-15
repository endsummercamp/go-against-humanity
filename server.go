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

	/* Custom Context */

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			db := utils.InitDB()
			cc := &utils.CustomContext{Context: c, Db: db, MatchManager: new(game.MatchManager)}
			return next(cc)
		}
	})



	e.Static("/public", "public")
	e.GET("/login", controllers.Login)
	e.POST("/login", controllers.DoLogin)
	e.GET("/signup", controllers.SignUp)
	e.POST("/signup", controllers.DoSignUp)
	e.GET("/logout", controllers.Logout)

	e.GET("/admin", controllers.Admin)
	e.GET("/matches", controllers.Matches)

	e.GET("/", controllers.Index)
	e.Logger.Debug(e.Start(":1323"))
}
