package ws

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/communication"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type client struct {
	id     uuid.UUID
	socket *websocket.Conn
	send   chan communication.Sender
	hub    HubInterface
	groups []uuid.UUID
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		// socket can read only communication message
		var msg communication.Message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		c.hub.Forward(&msg)
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := msg.Send(c.socket); err != nil {
			break
		}
	}
}
