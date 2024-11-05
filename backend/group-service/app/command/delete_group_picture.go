package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/google/uuid"
)

type DeleteGroupPictureCommand struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

type DeleteGroupPictureHandler struct {
	repo    database.GroupsRepository
	storage storage.Storage
}

func NewDeleteGroupPictureHandler(repo database.GroupsRepository, storage storage.Storage) DeleteGroupPictureHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if storage == nil {
		panic("storage is nil")
	}
	return DeleteGroupPictureHandler{repo: repo, storage: storage}
}

func (h DeleteGroupPictureHandler) Handle(ctx context.Context, cmd DeleteGroupPictureCommand) error {
	if err := h.repo.UpdateGroup(ctx, cmd.UserID, cmd.GroupID, func(g *models.Group) error {
		if !g.ChangePictureStateIfIncorrect(false) {
			return apperrors.NewBadRequest("group has no picture")
		}
		if err := h.storage.DeleteFile(ctx, cmd.GroupID.String()); err != nil {
			return apperrors.NewInternal(err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
