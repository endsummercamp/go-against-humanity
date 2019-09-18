package controllers

import (
	"github.com/ESCah/go-against-humanity/app/models"
	"log"
)

func (w *WebApp) GetLeaderboard(user_id int64) ([]models.User, error) {
	var board []models.User
	// Dirty hack to select the best 10 users, always including the current user
	_, err := w.Db.Select(&board, "SELECT * FROM (SELECT * FROM users ORDER BY user_id == ?, score DESC LIMIT 10) ORDER BY score DESC", string(user_id))
	if err != nil {
		log.Printf("%#v\n", err)
		return []models.User{}, err
	}

	return board, nil
}