package models

import (
	"sort"
	"sync"
	"time"
)

type Round struct {
	BlackCard      *BlackCard
	TimeFinishPick time.Time
	Wcs            map[*WhiteCard][]Juror
	Mutex	  	   sync.Mutex
	Voters			[]Juror

}

func (r *Round) AddCard(card *WhiteCard) bool {
	if _, ok := r.Wcs[card]; ok {
		return false
	}

	r.Wcs[card] = []Juror{}
	return true
}

func (r *Round) GetChoices() []*WhiteCard {
	var ret []*WhiteCard
	for card := range r.Wcs {
		ret = append(ret, card)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Text < ret[j].Text
	})
	return ret
}
