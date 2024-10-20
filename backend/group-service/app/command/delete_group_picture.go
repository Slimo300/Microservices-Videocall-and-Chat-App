package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
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
	emitter msgqueue.EventEmiter
}

func NewDeleteGroupPictureHandler(repo database.GroupsRepository, emitter msgqueue.EventEmiter, storage storage.Storage) DeleteGroupPictureHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	if storage == nil {
		panic("storage is nil")
	}
	return DeleteGroupPictureHandler{repo: repo, emitter: emitter, storage: storage}
}

func (h DeleteGroupPictureHandler) Handle(ctx context.Context, cmd DeleteGroupPictureCommand) error {
	if err := h.repo.UpdateGroup(ctx, cmd.GroupID, func(g *models.Group) error {
		member, ok := g.GetMemberByUserID(cmd.UserID)
		if !ok {
			return apperrors.NewNotFound("group not found")
		}
		if !member.CanUpdateGroup() {
			return apperrors.NewForbidden("user can't update group")
		}
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
