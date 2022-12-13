package user

import "errors"

var (
	ErrNoUser     = errors.New("no user found")
	ErrUserExists = errors.New("user already exists")
	ErrBadPass    = errors.New("invald password")
)

type UserMemoryRepository struct {
	data   map[string]*User
	lastID uint32
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: map[string]*User{
			"admin": {
				ID:       1,
				Login:    "admin",
				Type:     "admin",
				password: "admin",
			},
			"alexey": {
				ID:       2,
				Login:    "alexey",
				Type:     "user",
				password: "love",
			},
		},
		lastID: 2,
	}
}

func (repo *UserMemoryRepository) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}

	if u.password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}

func (repo *UserMemoryRepository) Register(login, pass string) (*User, error) {
	_, ok := repo.data[login]
	if ok {
		return nil, ErrUserExists
	}
	repo.lastID++
	u := &User{
		ID:       repo.lastID,
		Login:    login,
		password: pass,
		Type:     "user",
	}
	repo.data[login] = u

	return u, nil
}
