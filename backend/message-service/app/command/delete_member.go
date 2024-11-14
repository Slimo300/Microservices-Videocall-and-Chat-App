package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/google/uuid"
)

type DeleteMemberCommand struct {
	MemberID uuid.UUID
}

type DeleteMemberHandler struct {
	repo database.MessagesRepository
}

func NewDeleteMemberHandler(repo database.MessagesRepository) DeleteMemberHandler {
	if repo == nil {
		panic("nil repo")
	}
	return DeleteMemberHandler{repo: repo}
}

func (h *DeleteMemberHandler) Handle(ctx context.Context, cmd DeleteMemberCommand) error {
	return h.repo.DeleteMember(ctx, cmd.MemberID)
}
