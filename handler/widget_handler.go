package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

type WidgetHandler struct{}

func (*WidgetHandler) Alert(c echo.Context) error {
	content, err := os.ReadFile("public/alert.html")
	if err != nil {
		http.Error(c.Response().Writer, "Could not open requested file", http.StatusInternalServerError)
		return nil
	}

	fmt.Fprintf(c.Response().Writer, "%s", content)
	return nil
}

func (*WidgetHandler) Leaderboard(c echo.Context) error {
	
	content, err := os.ReadFile("public/leaderboard.html")
	if err != nil {
		http.Error(c.Response().Writer, "Could not open requested file", http.StatusInternalServerError)
		return nil
	}

	fmt.Fprintf(c.Response().Writer, "%s", content)
	return nil
}

func NewWidgetHandler() *WidgetHandler {
	return &WidgetHandler{}
}
