package user

type User struct {
	ID       uint32
	Login    string
	Type     string
	password string
}

type UserRepo interface {
	Authorize(login, pass string) (*User, error)
	Register(login, pass string) (*User, error)
}
