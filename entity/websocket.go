package entity

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

type WebSocketClient struct {
	*websocket.Conn
	donations chan Donation
}

func NewWebSocketClient(conn *websocket.Conn) WebSocketClient {
	return WebSocketClient{Conn: conn, donations: make(chan Donation)}
}

func (c *WebSocketClient) HandleIO() {
	var wg sync.WaitGroup
	wg.Add(2)

	go c.readPump(&wg)
	go c.writePump(&wg)

	wg.Wait()
}

func (c *WebSocketClient) SendDonation(donation Donation) {
	c.donations <- donation
}

func (c *WebSocketClient) Close() {
	close(c.donations)
	c.Conn.Close()
}

func (c *WebSocketClient) writePump(wg *sync.WaitGroup) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		wg.Done()
	}()

	for {
		select {
		case donation := <-c.donations:
			c.notify(donation)
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *WebSocketClient) readPump(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(s string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *WebSocketClient) notify(donation Donation) {
	err := c.Conn.WriteJSON(donation)
	if err != nil {
		log.Println(err)
	}
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*WebSocketClient]bool

	// Inbound messages from the clients.
	Broadcast chan Donation

	// Register requests from the clients.
	Register chan *WebSocketClient

	// Unregister requests from clients.
	Unregister chan *WebSocketClient
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Donation),
		Register:   make(chan *WebSocketClient),
		Unregister: make(chan *WebSocketClient),
		clients:    make(map[*WebSocketClient]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.donations)
			}
		case donation := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.donations <- donation:
				default:
					close(client.donations)
					delete(h.clients, client)
				}
			}
		}
	}
}
