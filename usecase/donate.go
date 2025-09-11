package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bickyeric/nyaweria/entity"
	liberr "github.com/bickyeric/nyaweria/errors"
	"github.com/bickyeric/nyaweria/repository"
	"github.com/google/uuid"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
)

type TopDonorsRequest struct {
	Username           string
	Limit              int
	StartTime, EndTime time.Time
}

type Donate interface {
	Donate(ctx context.Context, donation entity.Donation) error
	Summary(ctx context.Context, req TopDonorsRequest) ([]*entity.DonationSummary, error)
}

type donate struct {
	userRepo            repository.User
	donateRepo          repository.Donate
	speech              htgotts.Speech
	notificationUsecase Notification
}

func (u *donate) Summary(ctx context.Context, req TopDonorsRequest) ([]*entity.DonationSummary, error) {
	if req.Username == "" {
		return nil, liberr.UsernameEmptyErr
	}

	if req.Limit < 1 {
		req.Limit = 5
	}

	if req.StartTime.IsZero() {
		req.StartTime = time.Now().AddDate(0, 0, -1)
	}

	if req.EndTime.IsZero() {
		req.EndTime = time.Now()
	}

	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	summaries, err := u.donateRepo.Summary(ctx, repository.SummaryRequest{
		RecipientID: user.ID,
		Limit:       req.Limit,
		EndTime:     req.EndTime,
		StartTime:   req.StartTime,
	})
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func (u *donate) Donate(ctx context.Context, donation entity.Donation) error {
	slog.Info("donation received", slog.String("from", donation.From), slog.String("to", donation.To), slog.String("message", donation.Message))

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

func NewDonate(notificationUsecase Notification, userRepo repository.User, donateRepo repository.Donate, audioDirectory string) Donate {
	return &donate{
		userRepo:            userRepo,
		donateRepo:          donateRepo,
		speech:              htgotts.Speech{Folder: audioDirectory, Language: voices.Indonesian},
		notificationUsecase: notificationUsecase,
	}
}
