package communication

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const TIME_FORMAT = "2006-02-01 15:04:05"

// Message is a plain message in chat app
type Message struct {
	Group   uuid.UUID `json:"group"`
	Member  uuid.UUID `json:"member"`
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
