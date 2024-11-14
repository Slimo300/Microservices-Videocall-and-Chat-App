package query

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/google/uuid"
)

type GetPresignedGetRequestsQuery struct {
	UserID, GroupID uuid.UUID
	FileKeys        []storage.PresignGetFileInput
}

type GetPresignedGetRequestsHandler struct {
	repo    database.MessagesRepository
	storage storage.Storage
}

func NewGetPresignedGetRequestsHandler(repo database.MessagesRepository, storage storage.Storage) GetPresignedGetRequestsHandler {
	if repo == nil {
		panic("nil repo")
	}
	if storage == nil {
		panic("nil storage")
	}
	return GetPresignedGetRequestsHandler{
		repo:    repo,
		storage: storage,
	}
}

func (h *GetPresignedGetRequestsHandler) Handle(ctx context.Context, query GetPresignedGetRequestsQuery) ([]storage.PresignGetFileOutput, error) {
	if _, err := h.repo.GetUserGroupMember(ctx, query.UserID, query.GroupID); err != nil {
		return nil, err
	}
	return h.storage.GetPresignedGetRequests(ctx, query.FileKeys...)
}
