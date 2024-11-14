package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
)

type CreateMemberCommand struct {
	MemberID, UserID, GroupID uuid.UUID
	Username                  string
	Creator                   bool
}

type CreateMemberHandler struct {
	repo database.MessagesRepository
}

func NewCreateMemberHandler(repo database.MessagesRepository) CreateMemberHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CreateMemberHandler{repo: repo}
}

func (h *CreateMemberHandler) Handle(ctx context.Context, cmd CreateMemberCommand) error {
	return h.repo.CreateMember(ctx, models.NewMember(cmd.MemberID, cmd.UserID, cmd.GroupID, cmd.Username, cmd.Creator))
}
