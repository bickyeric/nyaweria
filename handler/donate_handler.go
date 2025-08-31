package handler

import (
	"net/http"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/bickyeric/nyaweria/usecase"
	"github.com/labstack/echo"
)

type DonateHandler struct {
	donateUsecase usecase.Donate
	userUsecase   usecase.User
}

func (h *DonateHandler) Index(c echo.Context) error {
	username := c.Param("streamer_username")

	user, err := h.userUsecase.GetByUsername(c.Request().Context(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.Render(http.StatusOK, "index.html", map[string]string{
		"streamer_username":    username,
		"profile_picture":      user.ProfilePicture,
		"streamer_name":        user.Name,
		"streamer_description": user.Description,
	})
}

func (h *DonateHandler) Donate(c echo.Context) error {
	var donation entity.Donation
	err := c.Bind(&donation)
	if err != nil {
		return err
	}

	err = h.donateUsecase.Donate(c.Request().Context(), donation)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, entity.ResponseBody{Message: "success giving donation"})
}

func (h *DonateHandler) Leaderboard(c echo.Context) error {
	summaries, err := h.donateUsecase.TopDonors(c.Request().Context(), c.FormValue("username"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, summaries)
}

func NewDonateHandler(donateUsecase usecase.Donate, userUsecase usecase.User) *DonateHandler {
	return &DonateHandler{
		donateUsecase: donateUsecase,
		userUsecase:   userUsecase,
	}
}
