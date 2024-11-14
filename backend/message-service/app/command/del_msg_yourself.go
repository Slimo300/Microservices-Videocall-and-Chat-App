package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/google/uuid"
)

type DeleteMessageForYourselfCommand struct {
	MessageID, UserID uuid.UUID
}

type DeleteMessageForYourselfHandler struct {
	repo database.MessagesRepository
}

func NewDeleteMessageForYourselfHandler(repo database.MessagesRepository) DeleteMessageForYourselfHandler {
	if repo == nil {
		panic("nil repo")
	}
	return DeleteMessageForYourselfHandler{repo: repo}
}

func (h *DeleteMessageForYourselfHandler) Handle(ctx context.Context, cmd DeleteMessageForYourselfCommand) error {
	return h.repo.DeleteMessageForYourself(ctx, cmd.UserID, cmd.MessageID)
}
