package models

type User struct {
	Id int64 `db:"user_id"`
	Username string
	PwHash string
	Admin bool
	UserType UserType
	Score int
}

type UserType int

const (
	PlayerType UserType = iota
	JurorType
)

func (u *User) IsAdmin() bool {
	return u.Admin
}