package ws

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const KEEP_ALIVE_INTERVAL = 50 * time.Second

type client struct {
	id     uuid.UUID
	socket *websocket.Conn
	send   chan Sender
	hub    WSHub
	groups map[uuid.UUID]struct{}
	ticker time.Ticker
}

// read reads messages received by socket
func (c *client) read() {
	defer c.socket.Close()
	for {
		// socket can read only communication message
		var msg Message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		c.hub.Forward(&msg)
	}
}

// write sends messages from server to clients
func (c *client) write() {
	defer c.socket.Close()
	for {
		select {
		case msg := <-c.send:
			if msg == nil {
				return
			}
			if err := msg.Send(c.socket); err != nil {
				log.Printf("Error when sending message through socket: %v\n", err)
			}
			c.ticker.Reset(1)
		case <-c.ticker.C:
			if err := c.socket.WriteMessage(websocket.PingMessage, []byte("KeepAlive")); err != nil {
				log.Printf("Error pinging client: %v\n", err)
			}
		}
	}
}
