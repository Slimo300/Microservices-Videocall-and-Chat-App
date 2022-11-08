package ws

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

// Sender is a type that is sent through websocket connection
type Sender interface {
	Send(*websocket.Conn) error
}

// Message is a plain message in chat app
type Message struct {
	Group   uuid.UUID `json:"group"`
	User    uuid.UUID `json:"user"`
	Message string    `json:"text"`
	Nick    string    `json:"nick"`
	When    string    `json:"created"`
}

// Send sends itself through websocket connection
func (m *Message) Send(ws *websocket.Conn) error {
	if err := ws.WriteJSON(m); err != nil {
		return err
	}
	return nil
}

func (m *Message) SetTime() {
	m.When = time.Now().Format(TIME_FORMAT)
}

type Action struct {
	ActionType string      `json:"type"`
	Payload    interface{} `json:"payload"`
}

func (a *Action) Send(ws *websocket.Conn) error {
	if err := ws.WriteJSON(a); err != nil {
		return err
	}
	return nil
}
