package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
)

type DeleteMemberCommand struct {
	UserID   uuid.UUID
	MemberID uuid.UUID
}

type DeleteMemberHandler struct {
	repo    database.GroupsRepository
	emitter msgqueue.EventEmiter
}

func NewDeleteMemberHandler(repo database.GroupsRepository, emitter msgqueue.EventEmiter) DeleteMemberHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return DeleteMemberHandler{repo: repo, emitter: emitter}
}

func (h DeleteMemberHandler) Handle(ctx context.Context, cmd DeleteMemberCommand) error {
	member, err := h.repo.GetMemberByID(ctx, cmd.UserID, cmd.MemberID)
	if err != nil {
		return apperrors.NewNotFound("member not found")
	}
	if err := h.repo.DeleteMember(ctx, cmd.UserID, cmd.MemberID); err != nil {
		return err
	}
	if err := h.emitter.Emit(events.MemberDeletedEvent{
		ID:      member.ID(),
		UserID:  member.UserID(),
		GroupID: member.GroupID(),
	}); err != nil {
		return err
	}
	return nil
}
