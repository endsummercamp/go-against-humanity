package models

import (
	"time"
	"sort"
)

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

func (r *Round) GetChoices() []*WhiteCard {
	var ret []*WhiteCard
	for card := range r.wcs {
		ret = append(ret, card)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Text < ret[j].Text
	})
	return ret
}