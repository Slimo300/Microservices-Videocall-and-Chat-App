package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
)

type SendInviteCommand struct {
	UserID   uuid.UUID
	GroupID  uuid.UUID
	TargetID uuid.UUID
}

type SendInviteHandler struct {
	repo    database.GroupsRepository
	emitter msgqueue.EventEmiter
}

func NewSendInviteHandler(repo database.GroupsRepository, emitter msgqueue.EventEmiter) SendInviteHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return SendInviteHandler{repo: repo, emitter: emitter}
}

func (h SendInviteHandler) Handle(ctx context.Context, cmd SendInviteCommand) error {
	invite := models.CreateInvite(cmd.UserID, cmd.TargetID, cmd.GroupID)
	invite, err := h.repo.CreateInvite(ctx, invite)
	if err != nil {
		return err
	}

	if err := h.emitter.Emit(events.InviteSentEvent{
		ID:       invite.ID(),
		IssuerID: invite.IssuerID(),
		Issuer: events.User{
			UserName:   invite.Issuer().Username(),
			HasPicture: invite.Issuer().HasPicture(),
		},
		TargetID: invite.TargetID(),
		GroupID:  invite.GroupID(),
		Group: events.Group{
			Name:       invite.Group().Name(),
			HasPicture: invite.Group().HasPicture(),
		},
		Status: int(invite.Status()),
	}); err != nil {
		return err
	}

	return nil
}
