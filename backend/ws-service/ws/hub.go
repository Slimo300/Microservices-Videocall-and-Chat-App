package ws

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

type WSHub struct {
	messageServerChan chan<- *Message
	forward           chan *Message
	join              chan *client
	leave             chan *client
	clients           map[*client]bool
}

func NewHub(messageChan chan<- *Message) *WSHub {
	return &WSHub{
		messageServerChan: messageChan,
		forward:           make(chan *Message),
		join:              make(chan *client),
		leave:             make(chan *client),
		clients:           make(map[*client]bool),
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.join:
			h.clients[client] = true
		case client := <-h.leave:
			delete(h.clients, client)
			close(client.send)
		case msg := <-h.forward:
			msg.SetTime()
			h.messageServerChan <- msg
			for client := range h.clients {
				for _, gr := range client.groups {
					if gr == msg.Group {
						client.send <- msg
					}
				}
			}
		}
	}
}

func (h *WSHub) Join(c *client) {
	h.join <- c
}

func (h *WSHub) Leave(c *client) {
	h.leave <- c
}
func (h *WSHub) Forward(msg *Message) {
	h.forward <- msg
}

func ServeWebSocket(w http.ResponseWriter, req *http.Request, h Hub, groups []uuid.UUID, id_user uuid.UUID) {

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}

	client := &client{
		id:     id_user,
		socket: socket,
		send:   make(chan Sender, messageBufferSize),
		hub:    h,
		groups: groups,
	}

	h.Join(client)
	defer func() { h.Leave(client) }()
	go client.write()
	client.read()
}
