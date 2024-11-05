package command

import (
	"context"
	"mime/multipart"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/google/uuid"
)

type SetGroupPicture struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
	File    multipart.File
}

type SetGroupPictureHandler struct {
	repo    database.GroupsRepository
	storage storage.Storage
}

func NewSetGroupPictureHandler(repo database.GroupsRepository, storage storage.Storage) SetGroupPictureHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if storage == nil {
		panic("storage is nil")
	}
	return SetGroupPictureHandler{repo: repo, storage: storage}
}

func (h SetGroupPictureHandler) Handle(ctx context.Context, cmd SetGroupPicture) error {
	if err := h.repo.UpdateGroup(ctx, cmd.UserID, cmd.GroupID, func(g *models.Group) error {
		if !g.ChangePictureStateIfIncorrect(true) {
			return nil
		}
		if err := h.storage.UploadFile(ctx, cmd.GroupID.String(), cmd.File, storage.PUBLIC_READ); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
