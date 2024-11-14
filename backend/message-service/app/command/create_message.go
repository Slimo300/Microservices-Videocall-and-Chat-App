package command

import (
	"context"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
)

type CreateMessageCommand struct {
	MessageID, MemberID, GroupID uuid.UUID
	Text, Nick                   string
	Posted                       time.Time
	Files                        []models.MessageFile
}

type CreateMessageHandler struct {
	repo database.MessagesRepository
}

func NewCreateMessageHandler(repo database.MessagesRepository) CreateMessageHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CreateMessageHandler{repo: repo}
}

func (h *CreateMessageHandler) Handle(ctx context.Context, cmd CreateMessageCommand) error {
	return h.repo.CreateMessage(ctx, models.NewMessage(cmd.MessageID, cmd.GroupID, cmd.MemberID, cmd.Text, cmd.Posted, cmd.Files))
}
