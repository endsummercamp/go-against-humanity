package models

import (
	"sort"
	"sync"
	"time"
)

type Round struct {
	BlackCard      *BlackCard
	TimeFinishPick time.Time
	Wcs            sync.Map // map[*WhiteCard][]*Juror
}

type WcsKey *WhiteCard
type WcsVal []*Juror

func (r *Round) AddCard(card *WhiteCard) bool {
	if _, ok := r.Wcs.Load(card); ok {
		return false
	}

	r.Wcs.Store(card, []*Juror{})
	return true
}

func (r *Round) GetChoices() []*WhiteCard {
	var ret []*WhiteCard
	r.Wcs.Range(func(_card, _ interface{}) bool {
		card := _card.(WcsKey)
		ret = append(ret, card)
		return true
	})
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Text < ret[j].Text
	})
	return ret
}