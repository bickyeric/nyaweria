package usecase

import (
	"context"
	"fmt"

	"github.com/bickyeric/nyaweria/entity"
)

type Donate interface {
	Donate(ctx context.Context, donation entity.Donation) error
}

type donate struct {
	notificationUsecase Notification
}

func (u *donate) Donate(ctx context.Context, donation entity.Donation) error {
	fmt.Printf("receive donation: %v\n", donation)

	return u.notificationUsecase.Send(ctx, donation)
}

func NewDonate(notificationUsecase Notification) Donate {
	return &donate{notificationUsecase: notificationUsecase}
}
