package database

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
)

type MessagesRepository interface {
	CreateMember(ctx context.Context, member models.Member) error
	UpdateMember(ctx context.Context, memberID uuid.UUID, updateFn func(m *models.Member) bool) error
	DeleteMember(ctx context.Context, memberID uuid.UUID) error
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error

	CreateMessage(ctx context.Context, message models.Message) error
	GetMessageByID(ctx context.Context, userID, messageID uuid.UUID) (models.Message, error)

	GetUserGroupMember(ctx context.Context, userID, groupID uuid.UUID) (models.Member, error)
	GetGroupMessages(ctx context.Context, userID, groupID uuid.UUID, offset, num int) ([]models.Message, error)
	DeleteMessageForYourself(ctx context.Context, userID, messageID uuid.UUID) error
	DeleteMessageForEveryone(ctx context.Context, userID, messageID uuid.UUID) error
}
