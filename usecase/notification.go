package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Notification interface {
	Add(ctx context.Context, ws *websocket.Conn, username string) (*entity.WebSocketPubsubClient, error)
	Delete(ctx context.Context, username string, conn *entity.WebSocketPubsubClient)
	Send(context.Context, entity.Donation) error
}

type notification struct {
	redisClient *redis.Client
	hubs        map[string]*entity.Hub
}

func (u *notification) Add(ctx context.Context, ws *websocket.Conn, username string) (*entity.WebSocketPubsubClient, error) {
	pubsubChannel := fmt.Sprintf("alert_%s", username)
	pubsub := u.redisClient.Subscribe(ctx, pubsubChannel)

	currentConn := entity.NewWebSocketPubsubClient(ws, pubsub)
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

func (u *notification) Delete(ctx context.Context, username string, conn *entity.WebSocketPubsubClient) {
	if hub, ok := u.hubs[username]; ok {
		hub.Unregister <- conn
	}
}

func (u *notification) Send(ctx context.Context, donation entity.Donation) error {
	jsonPayload, err := json.Marshal(donation)
	if err != nil {
		return err
	}

	pubsubChannel := fmt.Sprintf("alert_%s", donation.To)
	err = u.redisClient.Publish(ctx, pubsubChannel, string(jsonPayload)).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewNotification(redisClient *redis.Client) Notification {
	return &notification{
		hubs:        map[string]*entity.Hub{},
		redisClient: redisClient,
	}
}
