package usecase

import (
	"context"
	"fmt"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/google/uuid"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
)

type Donate interface {
	Donate(ctx context.Context, donation entity.Donation) error
}

type donate struct {
	speech              htgotts.Speech
	notificationUsecase Notification
}

func (u *donate) Donate(ctx context.Context, donation entity.Donation) error {
	fmt.Printf("receive donation: %v\n", donation)

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

func NewDonate(notificationUsecase Notification) Donate {
	return &donate{
		speech:              htgotts.Speech{Folder: "public/audio", Language: voices.Indonesian},
		notificationUsecase: notificationUsecase,
	}
}
