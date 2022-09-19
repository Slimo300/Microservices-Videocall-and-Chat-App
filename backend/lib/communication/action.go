package communication

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Sender interface decides which struct can be sent via websocket connection
type Sender interface {
	Send(*websocket.Conn) error
}

// Action represents type for signalizing changes to hub and further to frontend
type Action struct {
	Action string        `json:"action"` // pop or insert
	Group  uuid.UUID     `json:"group"`  // id_group
	User   uuid.UUID     `json:"-"`      // id_user
	Member models.Member `json:"member"` // member info for updates
	Invite models.Invite `json:"invite"` // invite
}

// Send sends itself through websocket connection
func (a *Action) Send(ws *websocket.Conn) error {
	if err := ws.WriteJSON(a); err != nil {
		return err
	}
	return nil
}
