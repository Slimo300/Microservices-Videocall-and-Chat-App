package query

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/google/uuid"
)

type GetPresignedPutRequestsQuery struct {
	UserID, GroupID uuid.UUID
	FileKeys        []storage.PresignPutFileInput
}

type GetPresignedPutRequestsHandler struct {
	repo    database.MessagesRepository
	storage storage.Storage
}

func NewGetPresignedPutRequestsHandler(repo database.MessagesRepository, storage storage.Storage) GetPresignedPutRequestsHandler {
	if repo == nil {
		panic("nil repo")
	}
	if storage == nil {
		panic("nil storage")
	}
	return GetPresignedPutRequestsHandler{
		repo:    repo,
		storage: storage,
	}
}

func (h *GetPresignedPutRequestsHandler) Handle(ctx context.Context, query GetPresignedPutRequestsQuery) ([]storage.PresignPutFileOutput, error) {
	if _, err := h.repo.GetUserGroupMember(ctx, query.UserID, query.GroupID); err != nil {
		return nil, err
	}
	return h.storage.GetPresignedPutRequests(ctx, query.FileKeys...)
}
