package utils

import (
	"github.com/ESCah/go-against-humanity/app/models"
)

func (c *CustomContext) GetUserByUsername(username string) *models.User {
	res, err := c.Db.Select(models.User{}, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		return nil
	}

	if res != nil && len(res) == 1 {
		return res[0].(*models.User)
	}

	return nil
}