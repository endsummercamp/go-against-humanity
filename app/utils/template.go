package utils

import (
	"html/template"
	"strings"
	"time"

	"github.com/ESCah/go-against-humanity/app/models"
)

var FuncMap = template.FuncMap{
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
