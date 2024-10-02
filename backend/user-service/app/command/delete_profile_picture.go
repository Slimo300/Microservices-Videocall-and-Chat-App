package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type DeleteProfilePicture struct {
	UserID uuid.UUID
}

type DeleteProfilePictureHandler struct {
	repo    database.UsersRepository
	storage storage.Storage
	emitter msgqueue.EventEmiter
}

func NewDeleteProfilePictureHandler(repo database.UsersRepository, storage storage.Storage, emitter msgqueue.EventEmiter) DeleteProfilePictureHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if storage == nil {
		panic("storage is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return DeleteProfilePictureHandler{repo: repo, storage: storage, emitter: emitter}
}

func (h DeleteProfilePictureHandler) Handle(ctx context.Context, cmd DeleteProfilePicture) error {
	if err := h.repo.UpdateUserByID(ctx, cmd.UserID, func(u *models.User) (*models.User, error) {
		if !u.ChangePictureStateIfIncorrect(false) {
			return nil, apperrors.NewForbidden("no picture to delete")
		}
		if err := h.storage.DeleteFile(ctx, cmd.UserID.String()); err != nil {
			return nil, err
		}
		if err := h.emitter.Emit(events.UserPictureModifiedEvent{ID: cmd.UserID, HasPicture: false}); err != nil {
			return nil, err
		}
		return u, nil
	}); err != nil {
		return err
	}
	return nil
}
