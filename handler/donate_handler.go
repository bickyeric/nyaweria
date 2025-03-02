package handler

import (
	"net/http"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/bickyeric/nyaweria/usecase"
	"github.com/labstack/echo"
)

type DonateHandler struct {
	donateUsecase usecase.Donate
}

func (*DonateHandler) Index(c echo.Context) error {
	streamerID := c.Param("streamer_id")

	return c.Render(http.StatusOK, "index.html", map[string]string{
		"streamer_id": streamerID,
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

func NewDonateHandler(donateUsecase usecase.Donate) *DonateHandler {
	return &DonateHandler{donateUsecase: donateUsecase}
}
