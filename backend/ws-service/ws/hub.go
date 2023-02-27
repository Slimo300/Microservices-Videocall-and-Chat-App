package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type WSHub struct {
	upgrader          *websocket.Upgrader
	actionServerChan  <-chan msgqueue.Event
	messageServerChan chan<- *Message
	forward           chan *Message
	join              chan *client
	leave             chan *client
	clients           map[*client]bool
}

func NewHub(messageChan chan<- *Message, actionChan <-chan msgqueue.Event, origin string) *WSHub {

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			if r.Header.Get("Origin") == origin {
				return true
			}
			return false
		}}

	return &WSHub{
		upgrader:          upgrader,
		messageServerChan: messageChan,
		actionServerChan:  actionChan,
		forward:           make(chan *Message),
		join:              make(chan *client),
		leave:             make(chan *client),
		clients:           make(map[*client]bool),
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case event := <-h.actionServerChan:
			switch e := event.(type) {
			case *events.GroupDeletedEvent:
				h.groupDeleted(*e)
			case *events.MemberCreatedEvent:
				h.memberAdded(*e)
			case *events.MemberDeletedEvent:
				h.memberDeleted(*e)
			case *events.MemberUpdatedEvent:
				h.memberUpdated(*e)
			case *events.InviteSentEvent:
				h.inviteSent(*e)
			case *events.InviteRespondedEvent:
				h.inviteResponded(*e)
			case *events.MessageDeletedEvent:
				h.messageDeleted(*e)
			default:
				log.Println("Unsupported Event Type: ", event.EventName())
			}
		case client := <-h.join:
			h.clients[client] = true
		case client := <-h.leave:
			delete(h.clients, client)
			close(client.send)
		case msg := <-h.forward:
			msg.Prepare()
			h.messageServerChan <- msg
			for client := range h.clients {
				if _, ok := client.groups[msg.Group]; ok {
					client.send <- msg
				}
				// for _, gr := range client.groups {
				// 	if gr == msg.Group {
				// 		client.send <- msg
				// 	}
				// }
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

func ServeWebSocket(w http.ResponseWriter, req *http.Request, h WSHub, groups []uuid.UUID, id_user uuid.UUID) {

	socket, err := h.upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}

	groupsMap := make(map[uuid.UUID]struct{})
	for _, group := range groups {
		groupsMap[group] = struct{}{}
	}

	client := &client{
		id:     id_user,
		socket: socket,
		send:   make(chan Sender, messageBufferSize),
		hub:    h,
		groups: groupsMap,
		ticker: *time.NewTicker(KEEP_ALIVE_INTERVAL),
	}

	h.Join(client)
	defer func() { h.Leave(client) }()
	go client.write()
	client.read()
}
