package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
)

type UpdateMemberCommand struct {
	MemberID                uuid.UUID
	DeletingMessages, Admin bool
}

type UpdateMemberHandler struct {
	repo database.MessagesRepository
}

func NewUpdateMemberHandler(repo database.MessagesRepository) UpdateMemberHandler {
	if repo == nil {
		panic("nil repo")
	}
	return UpdateMemberHandler{repo: repo}
}

func (h *UpdateMemberHandler) Handle(ctx context.Context, cmd UpdateMemberCommand) error {
	return h.repo.UpdateMember(ctx, cmd.MemberID, func(m *models.Member) bool {
		return m.UpdateRights(cmd.Admin, cmd.DeletingMessages)
	})
}
