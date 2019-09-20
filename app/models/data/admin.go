package data

import (
	"github.com/ESCah/go-against-humanity/app/models"
)

type AdminPageData struct {
	Header HeaderData
	Flash  FlashData
	User   models.User
}

type AdminUsersPageData struct {
	Header HeaderData
	Flash  FlashData
	User   models.User
	Users  []models.User
}
