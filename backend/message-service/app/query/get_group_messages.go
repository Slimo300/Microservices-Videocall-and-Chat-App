package query

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
)

type GetGroupMessagesQuery struct {
	UserID, GroupID uuid.UUID
	Num, Offset     int
}

type GetGroupMessagesHandler struct {
	repo database.MessagesRepository
}

func NewGetGroupMessagesHandler(repo database.MessagesRepository) GetGroupMessagesHandler {
	if repo == nil {
		panic("nil repo")
	}
	return GetGroupMessagesHandler{
		repo: repo,
	}
}

func (h *GetGroupMessagesHandler) Handle(ctx context.Context, query GetGroupMessagesQuery) ([]models.Message, error) {
	return h.repo.GetGroupMessages(ctx, query.UserID, query.GroupID, query.Offset, query.Num)
}
