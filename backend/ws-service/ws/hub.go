package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type Hub struct {
	serviceID uuid.UUID
	upgrader  *websocket.Upgrader
	emiter    msgqueue.EventEmiter
	eventChan <-chan msgqueue.Event
	forward   chan *Message
	join      chan *client
	leave     chan *client
	clients   map[*client]bool
}

func NewHub(eventChan <-chan msgqueue.Event, emiter msgqueue.EventEmiter, origin string) *Hub {

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == origin
		}}

	return &Hub{
		serviceID: uuid.New(),
		upgrader:  upgrader,
		eventChan: eventChan,
		emiter:    emiter,
		forward:   make(chan *Message),
		join:      make(chan *client),
		leave:     make(chan *client),
		clients:   make(map[*client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case event := <-h.eventChan:
			h.HandleEvent(event)
		case client := <-h.join:
			h.clients[client] = true
		case client := <-h.leave:
			delete(h.clients, client)
			close(client.send)
		case msg := <-h.forward:
			msg.ID = uuid.New()
			msg.When = time.Now()

			go func() {
				if err := h.EmitMessage(msg); err != nil {
					log.Printf("Error emiting message: %v", err)
				}
			}()

			for client := range h.clients {
				if _, ok := client.groups[msg.Member.GroupID]; ok {
					client.send <- msg
				}
			}
		}
	}
}

func (h *Hub) EmitMessage(msg *Message) error {

	var files []events.File
	for _, f := range msg.Files {
		files = append(files, events.File{Key: f.Key, Extension: f.Ext})
	}

	if err := h.emiter.Emit(events.MessageSentEvent{
		ID:        msg.ID,
		MemberID:  msg.Member.ID,
		GroupID:   msg.Member.GroupID,
		UserID:    msg.Member.UserID,
		Nick:      msg.Member.Username,
		Posted:    msg.When,
		Text:      msg.Message,
		Files:     files,
		ServiceID: h.serviceID,
	}); err != nil {
		return err
	}

	return nil
}

func (h *Hub) Join(c *client) {
	h.join <- c
}

func (h *Hub) Leave(c *client) {
	h.leave <- c
}
func (h *Hub) Forward(msg *Message) {
	h.forward <- msg
}

func (h *Hub) ServeWebSocket(w http.ResponseWriter, req *http.Request, groups []uuid.UUID, userID uuid.UUID) {

	socket, err := h.upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}

	groupsMap := make(map[uuid.UUID]bool)
	for _, group := range groups {
		groupsMap[group] = true
	}

	client := &client{
		id:     userID,
		socket: socket,
		send:   make(chan Sender, messageBufferSize),
		hub:    h,
		groups: groupsMap,
	}

	h.Join(client)
	defer func() { h.Leave(client) }()
	go client.write()
	client.read()
}
