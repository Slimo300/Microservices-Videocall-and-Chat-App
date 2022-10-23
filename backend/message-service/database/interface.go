package database

import (
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
	"github.com/google/uuid"
)

type DBLayer interface {
	GetGroupMessages(grouID uuid.UUID, offset, num int) ([]models.Message, error)
	AddMessage(memberID uuid.UUID, text string, when time.Time) error
	IsUserInGroup(userID, groupID uuid.UUID) bool
	//DeleteMessage
}
