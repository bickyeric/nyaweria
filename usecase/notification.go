package usecase

import (
	"context"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/gorilla/websocket"
)

type Notification interface {
	Add(ctx context.Context, ws *websocket.Conn, username string) (*entity.WebSocketClient, error)
	Delete(ctx context.Context, username string, conn *entity.WebSocketClient)
	Send(context.Context, entity.Donation) error
}

type notification struct {
	hubs map[string]*entity.Hub
}

func (u *notification) Add(ctx context.Context, ws *websocket.Conn, username string) (*entity.WebSocketClient, error) {
	currentConn := entity.NewWebSocketClient(ws)
	if hub, ok := u.hubs[username]; ok {
		hub.Register <- &currentConn
	} else {
		hub := entity.NewHub()
		go hub.Run()

		hub.Register <- &currentConn
		u.hubs[username] = hub
	}

	return &currentConn, nil
}

func (u *notification) Delete(ctx context.Context, username string, conn *entity.WebSocketClient) {
	if hub, ok := u.hubs[username]; ok {
		hub.Unregister <- conn
	}
}

func (u *notification) Send(ctx context.Context, donation entity.Donation) error {
	hub, ok := u.hubs[donation.To]
	if !ok {
		return nil
	}

	hub.Broadcast <- donation
	return nil
}

func NewNotification() Notification {
	return &notification{
		hubs: map[string]*entity.Hub{},
	}
}
