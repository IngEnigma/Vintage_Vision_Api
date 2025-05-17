package domain

type User struct {
	ID       uint
	Email    string
	Password string
	IsAdmin  bool

	Profiles []Profile
}
