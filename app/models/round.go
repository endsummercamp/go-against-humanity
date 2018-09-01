package models

type Round struct {
	bc *BlackCard
	wcs map[*WhiteCard][]*Juror
}

func (r *Round) AddCard(card *WhiteCard) bool {
	if _, ok := r.wcs[card]; ok {
		return false
	}

	r.wcs[card] = []*Juror{}
	return true
}