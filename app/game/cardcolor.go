//go:generate stringer -type=CardColor
package game

type CardColor int

const (
	WHITE_CARD CardColor = iota
	BLACK_CARD
)
