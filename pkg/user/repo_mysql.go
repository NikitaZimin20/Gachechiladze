package user

import (
	"database/sql"
	"errors"
)

var (
	ErrNoUser     = errors.New("no user found")
	ErrBadPass    = errors.New("invald password")
	ErrUserExists = errors.New("user exists")
)

type UserMysqlRepository struct {
	DB *sql.DB
}

func NewMysqlRepo(db *sql.DB) *UserMysqlRepository {
	return &UserMysqlRepository{DB: db}
}

func (repo *UserMysqlRepository) Authorize(login, pass string) (*User, error) {
	user := &User{}
	if err := repo.DB.QueryRow("SELECT id, username, password, type FROM users WHERE username = ?", login).
		Scan(&user.ID, &user.Login, &user.password, &user.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoUser
		}
		return nil, err
	}
	if user.password != pass {
		return nil, ErrBadPass
	}

	return user, nil
}

func (repo *UserMysqlRepository) Register(login, pass string) (*User, error) {
	result, err := repo.DB.Exec(
		"INSERT INTO users (username, password, type) VALUES (?, ?, ?)",
		login,
		pass,
		"user",
	)
	if err != nil {
		return &User{}, ErrUserExists
	}

	res, err := result.LastInsertId()
	if err != nil {
		return &User{}, err
	}

	user := &User{
		ID:       uint32(res),
		Login:    login,
		password: pass,
		Type:     "user",
	}

	return user, nil
}
