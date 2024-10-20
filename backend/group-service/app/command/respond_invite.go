package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
)

type RespondInvite struct {
	UserID   uuid.UUID
	InviteID uuid.UUID
	Answer   bool
}

type RespondInviteHandler struct {
	repo    database.GroupsRepository
	emitter msgqueue.EventEmiter
}

func NewRespondInviteHandler(repo database.GroupsRepository, emitter msgqueue.EventEmiter) RespondInviteHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return RespondInviteHandler{repo: repo, emitter: emitter}
}

func (h RespondInviteHandler) Handle(ctx context.Context, cmd RespondInvite) error {
	var invite models.Invite
	var member models.Member
	if err := h.repo.UpdateInvite(ctx, cmd.InviteID, func(i *models.Invite) (*models.Member, error) {
		if err := i.AnswerInvite(cmd.UserID, cmd.Answer); err != nil {
			return nil, err
		}
		invite = *i
		member = models.NewMember(cmd.UserID, i.GroupID())
		return &member, nil
	}); err != nil {
		return err
	}
	if err := h.emitter.Emit(events.MemberCreatedEvent{
		ID:      member.ID(),
		GroupID: member.GroupID(),
		UserID:  member.UserID(),
		User: events.User{
			UserName:   member.User().Username(),
			HasPicture: member.User().HasPicture(),
		},
		Creator: member.Creator(),
	}); err != nil {
		return err
	}
	if err := h.emitter.Emit(events.InviteRespondedEvent{
		ID:       invite.ID(),
		IssuerID: invite.IssuerID(),
		TargetID: invite.TargetID(),
		Target: events.User{
			UserName:   invite.Target().Username(),
			HasPicture: invite.Target().HasPicture(),
		},
		GroupID: invite.GroupID(),
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
