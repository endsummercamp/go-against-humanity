package models

type Player struct {
	User   *User
	Cards  []*WhiteCard
}
