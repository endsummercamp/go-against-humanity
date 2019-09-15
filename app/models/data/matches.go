package data

import (
	"github.com/ESCah/go-against-humanity/app/models"
)

type MatchesPageData struct {
	Header HeaderData
	Flash  FlashData
	User   models.User
	Matches []*models.Match
}
