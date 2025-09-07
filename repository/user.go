package repository

import (
	"context"
	"database/sql"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/doug-martin/goqu/v9"
)

//go:generate mockgen -source=user.go -destination=mock/user.go
type User interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type user struct {
	db *sql.DB
}

func (u *user) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query, args, err := goqu.From("users").
		Select("id", "username", "name", "profile_picture", "description").
		Where(goqu.C("username").Eq(username)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var user entity.User

	row := u.db.QueryRow(query, args...)
	if err = row.Scan(&user.ID, &user.Username, &user.Name, &user.ProfilePicture, &user.Description); err != nil {
		return nil, err
	}

	return &user, nil
}

func NewUser(db *sql.DB) User {
	return &user{
		db: db,
	}
}
