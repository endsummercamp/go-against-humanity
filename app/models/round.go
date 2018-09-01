package models

type Round struct {
	bc *BlackCard
	wcs map[*WhiteCard][]*Juror
}