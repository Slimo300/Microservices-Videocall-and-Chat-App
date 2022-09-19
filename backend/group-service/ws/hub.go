package ws

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/communication"
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

type Hub struct {
	actionServerChan  <-chan *communication.Action
	messageServerChan chan<- *communication.Message
	forward           chan *communication.Message
	join              chan *client
	leave             chan *client
	clients           map[*client]bool
}

func NewHub(messageChan chan<- *communication.Message, actionChan <-chan *communication.Action) *Hub {
	return &Hub{
		actionServerChan:  actionChan,
		messageServerChan: messageChan,
		forward:           make(chan *communication.Message),
		join:              make(chan *client),
		leave:             make(chan *client),
		clients:           make(map[*client]bool),
	}
}

func (h *Hub) Run() {
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
		case msg := <-h.actionServerChan:
			switch msg.Action {
			case "DELETE_GROUP":
				h.GroupDeleted(msg.Group)
			case "CREATE_GROUP":
				h.GroupCreated(msg.User, msg.Group)
			case "ADD_MEMBER":
				h.MemberAdded(msg.Member)
			case "DELETE_MEMBER":
				h.MemberDeleted(msg.Member)
			case "SEND_INVITE":
				h.SendGroupInvite(msg.Invite)
			}
		}
	}
}

func (h *Hub) Join(c *client) {
	h.join <- c
}

func (h *Hub) Leave(c *client) {
	h.leave <- c
}
func (h *Hub) Forward(msg *communication.Message) {
	h.forward <- msg
}

func ServeWebSocket(w http.ResponseWriter, req *http.Request, h HubInterface, groups []uuid.UUID, id_user uuid.UUID) {

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}

	client := &client{
		id:     id_user,
		socket: socket,
		send:   make(chan communication.Sender, messageBufferSize),
		hub:    h,
		groups: groups,
	}

	h.Join(client)
	defer func() { h.Leave(client) }()
	go client.write()
	client.read()
}
