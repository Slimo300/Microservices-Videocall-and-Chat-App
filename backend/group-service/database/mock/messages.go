package mock

import (
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
)

func (m *MockDB) GetGroupMessages(groupID uuid.UUID, offset, num int) (messages []models.Message, err error) {
	for _, message := range m.Messages {
		for _, member := range m.Members {
			if message.MemberID == member.ID && member.GroupID == groupID {
				message.Member = member
				messages = append(messages, message)
			}
		}
	}
	return messages, nil
}

func (m *MockDB) AddMessage(memberID uuid.UUID, text string, when time.Time) error {
	m.Messages = append(m.Messages, models.Message{
		ID:       uuid.New(),
		Posted:   when,
		Text:     text,
		MemberID: memberID,
	})

	return nil
}
