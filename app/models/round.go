package models

type Round struct {
	bc *BlackCard
	wcs map[*WhiteCard][]*Juror
}

func (r *Round) AddCard(card *WhiteCard){
	if _, ok := r.wcs[card]; ok {
		return
	}

	r.wcs[card] = []*Juror{}
}