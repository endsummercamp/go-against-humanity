package models

type Player struct {
	User		*User
	Points     	int
	Cards      	[]*WhiteCard
}