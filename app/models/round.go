package models

import "time"

type Round struct {
	BlackCard *BlackCard
	TimeFinishPick time.Time
	wcs       map[*WhiteCard][]*Juror
}

func (r *Round) AddCard(card *WhiteCard) bool {
	if _, ok := r.wcs[card]; ok {
		return false
	}

	r.wcs[card] = []*Juror{}
	return true
}