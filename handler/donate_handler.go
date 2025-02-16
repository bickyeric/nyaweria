package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/bickyeric/nyaweria/usecase"
	"github.com/labstack/echo"
)

type DonateHandler struct {
	donateUsecase usecase.Donate
}

func (*DonateHandler) Index(c echo.Context) error {
	content, err := os.ReadFile("public/index.html")
	if err != nil {
		http.Error(c.Response().Writer, "Could not open requested file", http.StatusInternalServerError)
		return nil
	}

	fmt.Fprintf(c.Response().Writer, "%s", content)
	return nil
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
