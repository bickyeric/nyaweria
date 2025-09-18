package usecase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/bickyeric/nyaweria/repository"
	mock_repository "github.com/bickyeric/nyaweria/repository/mock"
	"github.com/bickyeric/nyaweria/usecase"
	mock_usecase "github.com/bickyeric/nyaweria/usecase/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type DonateSuite struct {
	suite.Suite
}

func (s *DonateSuite) TestSummaryEmptyUsername() {
	ctrl := gomock.NewController(s.T())
	audioDir := os.TempDir()

	mockNotification := mock_usecase.NewMockNotification(ctrl)
	mockUserRepo := mock_repository.NewMockUser(ctrl)
	mockDonateRepo := mock_repository.NewMockDonate(ctrl)

	donateUsecase := usecase.NewDonate(mockNotification, mockUserRepo, mockDonateRepo, audioDir)
	_, err := donateUsecase.Summary(context.Background(), usecase.TopDonorsRequest{})
	s.ErrorContains(err, "username cannot be empty")
}

func (s *DonateSuite) TestSummarySuccess() {
	ctrl := gomock.NewController(s.T())
	audioDir := os.TempDir()

	mockNotification := mock_usecase.NewMockNotification(ctrl)
	mockUserRepo := mock_repository.NewMockUser(ctrl)
	mockDonateRepo := mock_repository.NewMockDonate(ctrl)

	timeNow := time.Now()

	mockUserRepo.EXPECT().GetByUsername(context.Background(), "bickyeric").Return(&entity.User{ID: "123"}, nil)
	mockDonateRepo.EXPECT().Summary(context.Background(), repository.SummaryRequest{
		RecipientID: "123",
		Limit:       5,
		StartTime:   timeNow.AddDate(0, 0, -1),
		EndTime:     timeNow,
	}).Return([]*entity.DonationSummary{
		{
			Sender: "user 1",
			Sum:    250000,
		},
	}, nil)

	donateUsecase := usecase.NewDonate(mockNotification, mockUserRepo, mockDonateRepo, audioDir)
	summaries, err := donateUsecase.Summary(context.Background(), usecase.TopDonorsRequest{
		Username:  "bickyeric",
		Limit:     5,
		StartTime: timeNow.AddDate(0, 0, -1),
		EndTime:   timeNow,
	})
	s.NoError(err)
	s.Len(summaries, 1)
	s.Equal(&entity.DonationSummary{Sender: "user 1", Sum: 250000}, summaries[0])
}

func (s *DonateSuite) TestDonateSuccess() {
	ctrl := gomock.NewController(s.T())
	audioDir := os.TempDir()

	mockNotification := mock_usecase.NewMockNotification(ctrl)
	mockUserRepo := mock_repository.NewMockUser(ctrl)
	mockDonateRepo := mock_repository.NewMockDonate(ctrl)

	mockDonateRepo.EXPECT().Create(context.Background(), gomock.Any()).Return(nil)
	mockUserRepo.EXPECT().GetByUsername(context.Background(), "bickyeric").Return(&entity.User{ID: "123"}, nil)
	mockNotification.EXPECT().Send(context.Background(), gomock.Any()).Return(nil)

	donateUsecase := usecase.NewDonate(mockNotification, mockUserRepo, mockDonateRepo, audioDir)
	err := donateUsecase.Donate(context.Background(), entity.Donation{
		From:    "user 2",
		To:      "bickyeric",
		Amount:  "125000",
		Message: "halo testing bang",
	})
	s.Error(err)
}

func TestDonateSuite(t *testing.T) {
	suite.Run(t, &DonateSuite{})
}
