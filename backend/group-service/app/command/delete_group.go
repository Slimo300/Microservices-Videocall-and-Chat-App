package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
)

type DeleteGroupCommand struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

type DeleteGroupHandler struct {
	repo    database.GroupsRepository
	emitter msgqueue.EventEmiter
}

func NewDeleteGroupHandler(repo database.GroupsRepository, emitter msgqueue.EventEmiter) DeleteGroupHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return DeleteGroupHandler{repo: repo, emitter: emitter}
}

func (h DeleteGroupHandler) Handle(ctx context.Context, cmd DeleteGroupCommand) error {
	if err := h.repo.DeleteGroup(ctx, cmd.UserID, cmd.GroupID); err != nil {
		return err
	}
	if err := h.emitter.Emit(events.GroupDeletedEvent{
		ID: cmd.GroupID,
	}); err != nil {
		return apperrors.NewInternal(err)
	}
	return nil
}
