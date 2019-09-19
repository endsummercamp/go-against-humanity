package utils

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func GetUsername(c echo.Context) string {
	s, _ := session.Get("session", c)
	if x, ok := s.Values["user"]; ok {
		return x.(string)
	}

	return ""
}

func IsLoggedIn(c echo.Context) bool {
	return GetUsername(c) != ""
}
