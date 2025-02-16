package entity

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	*websocket.Conn
	donations chan Donation
}

func NewWebSocketConnection(conn *websocket.Conn) WebSocketConnection {
	return WebSocketConnection{Conn: conn, donations: make(chan Donation)}
}

func (c *WebSocketConnection) HandleIO() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	for {
		fmt.Println("waiting for donation")
		donation, ok := <-c.donations
		if !ok {
			fmt.Println("connection not okay")
		}

		fmt.Println("sending notification")
		c.notify(donation)
	}
}

func (c *WebSocketConnection) SendDonation(donation Donation) {
	c.donations <- donation
}

func (c *WebSocketConnection) notify(donation Donation) {
	err := c.Conn.WriteJSON(donation)
	if err != nil {
		log.Println(err)
	}
}
