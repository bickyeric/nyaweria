package entity

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 1 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type WebSocketPubsubClient struct {
	websocketConnection *websocket.Conn
	pubsubConnection    *redis.PubSub
}

func NewWebSocketPubsubClient(websocketConnection *websocket.Conn, pubsubConnection *redis.PubSub) WebSocketPubsubClient {
	return WebSocketPubsubClient{websocketConnection: websocketConnection, pubsubConnection: pubsubConnection}
}

func (c *WebSocketPubsubClient) HandleIO() {
	var wg sync.WaitGroup
	wg.Add(3)

	go c.readPump(&wg)
	go c.writePump(&wg)
	go c.readPubSub(&wg)

	wg.Wait()
}

func (c *WebSocketPubsubClient) Close() {
	slog.Info("closing websocket and pubsub connection")
	if err := c.websocketConnection.Close(); err != nil {
		slog.Error("error closing websocket", slog.String("error", err.Error()))
	}

	if err := c.pubsubConnection.Close(); err != nil {
		slog.Error("error closing pubsub", slog.String("error", err.Error()))
	}
}

func (c *WebSocketPubsubClient) readPubSub(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	slog.Info("pubsub connection opened")

	for msg := range c.pubsubConnection.Channel() {
		var donation Donation

		err := json.Unmarshal([]byte(msg.Payload), &donation)
		if err != nil {
			slog.Error("error unmarshal payload", slog.String("error", err.Error()))
			continue
		}

		c.notify(donation)
	}

	slog.Info("pubsub connection closed")
}

func (c *WebSocketPubsubClient) writePump(wg *sync.WaitGroup) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		wg.Done()
	}()

	for range ticker.C {
		if err := c.websocketConnection.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
			slog.Error("error writePump", slog.String("error", err.Error()))
			c.Close()
			return
		}
		if err := c.websocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
			slog.Error("error writePump", slog.String("error", err.Error()))
			c.Close()
			return
		}
	}
}

func (c *WebSocketPubsubClient) readPump(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	c.websocketConnection.SetReadLimit(maxMessageSize)
	if err := c.websocketConnection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		slog.Error("error readPump", slog.String("error", err.Error()))
		c.Close()
		return
	}
	c.websocketConnection.SetPongHandler(func(s string) error {
		return c.websocketConnection.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, _, err := c.websocketConnection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("failed read websocket message", slog.String("error", err.Error()))
			}
			c.Close()
			break
		}
	}
}

func (c *WebSocketPubsubClient) notify(donation Donation) {
	err := c.websocketConnection.WriteJSON(donation)
	if err != nil {
		slog.Error("error notify websocket", slog.String("error", err.Error()))
	}
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*WebSocketPubsubClient]bool

	// Register requests from the clients.
	Register chan *WebSocketPubsubClient

	// Unregister requests from clients.
	Unregister chan *WebSocketPubsubClient
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *WebSocketPubsubClient),
		Unregister: make(chan *WebSocketPubsubClient),
		clients:    make(map[*WebSocketPubsubClient]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			slog.Info("registering")
			h.clients[client] = true
		case client := <-h.Unregister:
			slog.Info("unregistering")
			delete(h.clients, client)
		}
	}
}
