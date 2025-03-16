package usecase

import (
	"context"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/bickyeric/nyaweria/repository"
)

type User interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type user struct {
	userRepository repository.User
}

func (u *user) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return u.userRepository.GetByUsername(ctx, username)
}

func NewUser(userRepository repository.User) User {
	return &user{userRepository}
}
