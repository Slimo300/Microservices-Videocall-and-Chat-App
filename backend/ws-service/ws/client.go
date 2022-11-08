package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type client struct {
	id     uuid.UUID
	socket *websocket.Conn
	send   chan Sender
	hub    Hub
	groups []uuid.UUID
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
	for msg := range c.send {
		if err := msg.Send(c.socket); err != nil {
			break
		}
	}
}
