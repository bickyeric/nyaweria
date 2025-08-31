package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/bickyeric/nyaweria/repository"
	"github.com/google/uuid"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
)

type Donate interface {
	Donate(ctx context.Context, donation entity.Donation) error
	TopDonors(ctx context.Context, username string) ([]*entity.DonationSummary, error)
}

type donate struct {
	userRepo            repository.User
	donateRepo          repository.Donate
	speech              htgotts.Speech
	notificationUsecase Notification
}

func (u *donate) TopDonors(ctx context.Context, username string) ([]*entity.DonationSummary, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	summaries, err := u.donateRepo.Summary(ctx, repository.SummaryRequest{
		RecipientID: user.ID,
		Limit:       5,
		EndTime:     time.Now(),
		StartTime:   time.Now().Add(-24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func (u *donate) Donate(ctx context.Context, donation entity.Donation) error {
	fmt.Printf("receive donation: %v\n", donation)

	recipient, err := u.userRepo.GetByUsername(ctx, donation.To)
	if err != nil {
		return err
	}
	donation.RecipientID = recipient.ID

	err = u.donateRepo.Create(ctx, &donation)
	if err != nil {
		return err
	}

	filename := uuid.New().String()

	giver := "Seseorang"
	if donation.From != "" {
		giver = donation.From
	}

	audioPath, err := u.speech.CreateSpeechFile(fmt.Sprintf("%s baru saja memberikan %s. %s", giver, donation.Amount, donation.Message), filename)
	if err != nil {
		return err
	}
	donation.AudioPath = audioPath

	return u.notificationUsecase.Send(ctx, donation)
}

func NewDonate(notificationUsecase Notification, userRepo repository.User, donateRepo repository.Donate) Donate {
	return &donate{
		userRepo:            userRepo,
		donateRepo:          donateRepo,
		speech:              htgotts.Speech{Folder: "public/audio", Language: voices.Indonesian},
		notificationUsecase: notificationUsecase,
	}
}
