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

	var mm = new(game.MatchManager)
	var ws = controllers.MakeSocketServer(mm)

	w := controllers.WebApp{
		MatchManager: mm,
		Echo:         e,
		Db:           utils.InitDB(),
		Ws:           ws,
	}

	e.Static("/public", "public")
	e.GET("/login", w.Login)
	e.POST("/login", w.DoLogin)
	e.GET("/signup", w.SignUp)
	e.POST("/signup", w.DoSignUp)
	e.GET("/logout", w.Logout)

	e.GET("/admin", w.Admin)
	e.GET("/admin/users", w.AdminUsers)
	// Todo: make RESTful
	e.GET("/admin/matches/new", w.AdminNewMatch)
	e.PUT("/admin/matches/:id/new_black_card", w.NewBlackCard)
	e.PUT("/admin/matches/:id/end_voting", w.EndVoting)

	e.GET("/matches", w.Matches)
	e.GET("/matches/join/:id", w.JoinMatch)
	e.GET("/matches/join_latest", w.JoinLatestMatch)

	e.GET("/mycards", w.MatchCards)
	e.PUT("/matches/:match_id/pick_card/:card_id", w.PickCard)
	e.PUT("/matches/:match_id/vote_card/:card_id", w.VoteCard)
	e.GET("/", w.Index)
	// The ws start is not blocking, but the echo start is, so ws goes first
	ws.Start()

	e.Logger.Debug(e.Start(":1323"))
}
