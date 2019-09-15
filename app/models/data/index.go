package data

import (
	"github.com/ESCah/go-against-humanity/app/models"
)

type IndexPageData struct {
	Header HeaderData
	Flash  FlashData
	User   models.User
}
