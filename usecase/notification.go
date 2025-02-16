package usecase

import (
	"context"
	"fmt"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/gorilla/websocket"
)

type Notification interface {
	Add(ctx context.Context, ws *websocket.Conn, username string) (entity.WebSocketConnection, error)
	Send(context.Context, entity.Donation) error
}

type notification struct {
	conns map[string]entity.WebSocketConnection
}

func (u *notification) Add(ctx context.Context, ws *websocket.Conn, username string) (entity.WebSocketConnection, error) {
	currentConn := entity.NewWebSocketConnection(ws)
	u.conns[username] = currentConn

	return currentConn, nil
}

func (u *notification) Send(ctx context.Context, donation entity.Donation) error {
	con, ok := u.conns[donation.To]
	if !ok {
		fmt.Println("streamer connection not found!!!")
	}

	con.SendDonation(donation)
	return nil
}

func NewNotification() Notification {
	return &notification{
		conns: map[string]entity.WebSocketConnection{},
	}
}
