//go:generate stringer -type=CardColor
package models

type CardColor int

const (
	WHITE_CARD CardColor = iota
	BLACK_CARD
)
