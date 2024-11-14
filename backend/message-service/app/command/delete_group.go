package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/google/uuid"
)

type DeleteGroupCommand struct {
	GroupID uuid.UUID
}

type DeleteGroupHandler struct {
	repo    database.MessagesRepository
	storage storage.Storage
}

func NewDeleteGroupHandler(repo database.MessagesRepository, storage storage.Storage) DeleteGroupHandler {
	if repo == nil {
		panic("nil repo")
	}
	if storage == nil {
		panic("nil storage")
	}
	return DeleteGroupHandler{repo: repo, storage: storage}
}

func (h *DeleteGroupHandler) Handle(ctx context.Context, cmd DeleteGroupCommand) error {
	if err := h.storage.DeleteFilesByPrefix(ctx, cmd.GroupID.String()+"/"); err != nil {
		return err
	}
	return h.repo.DeleteGroup(ctx, cmd.GroupID)
}
