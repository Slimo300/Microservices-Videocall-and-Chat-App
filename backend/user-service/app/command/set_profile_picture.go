package command

import (
	"context"
	"mime/multipart"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type SetProfilePicture struct {
	UserID uuid.UUID
	File   multipart.File
}

type SetProfilePictureHandler struct {
	repo    database.UsersRepository
	storage storage.Storage
	emitter msgqueue.EventEmiter
}

func NewSetProfilePictureHandler(repo database.UsersRepository, storage storage.Storage, emitter msgqueue.EventEmiter) SetProfilePictureHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if storage == nil {
		panic("storage is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return SetProfilePictureHandler{repo: repo, storage: storage, emitter: emitter}
}

func (h SetProfilePictureHandler) Handle(ctx context.Context, cmd SetProfilePicture) error {
	if err := h.repo.UpdateUserByID(ctx, cmd.UserID, func(u *models.User) (*models.User, error) {
		// Here we check whether user state was changed at all, if it was we emit event and continue to return user to be updated
		// if it wasn't we return nil, nil to signalize that there is no need to perform update in database
		if u.ChangePictureStateIfIncorrect(true) {
			if err := h.emitter.Emit(events.UserPictureModifiedEvent{ID: cmd.UserID, HasPicture: true}); err != nil {
				return nil, err
			}
		} else {
			return nil, nil
		}
		if err := h.storage.UploadFile(ctx, cmd.UserID.String(), cmd.File, storage.PUBLIC_READ); err != nil {
			return nil, err
		}
		return u, nil
	}); err != nil {
		return err
	}
	return nil
}
