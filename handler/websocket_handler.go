package handler

import (
	"log/slog"

	"github.com/bickyeric/nyaweria/usecase"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var upgrader = websocket.Upgrader{}

type WebsocketHandler struct {
	notification usecase.Notification
}

func (h *WebsocketHandler) Handle(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := ws.Close(); err != nil {
			slog.Error("error closing websocket", slog.String("error", err.Error()))
		}
	}()

	username := c.Request().URL.Query().Get("username")
	websocketConnection, err := h.notification.Add(c.Request().Context(), ws, username)
	if err != nil {
		return err
	}

	// block the request
	websocketConnection.HandleIO()

	h.notification.Delete(c.Request().Context(), username, websocketConnection)

	return nil
}

func NewWebsocketHandler(notification usecase.Notification) *WebsocketHandler {
	return &WebsocketHandler{notification: notification}
}
