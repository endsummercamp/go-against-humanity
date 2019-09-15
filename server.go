package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/ESCah/go-against-humanity/app/controllers"
	"github.com/ESCah/go-against-humanity/app/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
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

	e.Static("/public", "public")
	e.GET("/login", controllers.Login)
	e.POST("/login", controllers.DoLogin)
	e.GET("/signup", controllers.SignUp)
	e.POST("/signup", controllers.DoSignUp)
	e.GET("/logout", controllers.Logout)

	e.GET("/", controllers.Index)
	e.Logger.Debug(e.Start(":1323"))
}
