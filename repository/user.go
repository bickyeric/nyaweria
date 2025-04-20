package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bickyeric/nyaweria/entity"
)

type User interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type user struct {
	users map[string]*entity.User
	db    *sql.DB
}

func (u *user) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, ok := u.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func NewUser(db *sql.DB) User {
	return &user{
		db: db,
		users: map[string]*entity.User{
			"bickyeric": {
				Username:       "bickyeric",
				Name:           "Bicky Eric Kantona",
				ProfilePicture: "https://flowbite.com/docs/images/people/profile-picture-3.jpg",
				Description:    "Programmer Magang",
			},
			"streamertesting": {
				Username:       "streamertesting",
				Name:           "Streamer Testing",
				ProfilePicture: "https://img.icons8.com/?size=200&id=44442&format=png&color=000000",
				Description:    "Programmer Golang",
			},
		},
	}
}
