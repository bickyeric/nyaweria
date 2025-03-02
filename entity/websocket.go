package entity

import (
	"fmt"
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

type WebSocketConnection struct {
	*websocket.Conn
	donations chan Donation
}

func NewWebSocketConnection(conn *websocket.Conn) WebSocketConnection {
	return WebSocketConnection{Conn: conn, donations: make(chan Donation)}
}

func (c *WebSocketConnection) HandleIO() {
	var wg sync.WaitGroup
	wg.Add(2)

	go c.readPump(&wg)
	go c.writePump(&wg)

	wg.Wait()
}

func (c *WebSocketConnection) SendDonation(donation Donation) {
	c.donations <- donation
}

func (c *WebSocketConnection) Close() {
	close(c.donations)
	c.Conn.Close()
}

func (c *WebSocketConnection) writePump(wg *sync.WaitGroup) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		wg.Done()
	}()

	fmt.Println("waiting for donation")
	for {
		select {
		case donation, ok := <-c.donations:
			if !ok {
				fmt.Println("connection not okay")
			}

			fmt.Println("sending notification")
			c.notify(donation)
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("error sending ping message", err)
				return
			}
		}
	}
}

func (c *WebSocketConnection) readPump(wg *sync.WaitGroup) {
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

func (c *WebSocketConnection) notify(donation Donation) {
	err := c.Conn.WriteJSON(donation)
	if err != nil {
		log.Println(err)
	}
}
