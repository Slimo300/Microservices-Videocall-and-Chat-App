package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
)

type GrantRights struct {
	UserID           uuid.UUID
	MemberID         uuid.UUID
	Adding           bool
	DeletingMembers  bool
	DeletingMessages bool
	Muting           bool
	Admin            bool
}

type GrantRightsHandler struct {
	repo    database.GroupsRepository
	emitter msgqueue.EventEmiter
}

func NewGrantRightsHandler(repo database.GroupsRepository, emitter msgqueue.EventEmiter) GrantRightsHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return GrantRightsHandler{repo: repo, emitter: emitter}
}

func (h GrantRightsHandler) Handle(ctx context.Context, cmd GrantRights) error {
	var member *models.Member
	if err := h.repo.UpdateMember(ctx, cmd.UserID, cmd.MemberID, func(m *models.Member) error {
		m.ApplyRights(models.MemberRights{
			Adding:           cmd.Adding,
			DeletingMessages: cmd.DeletingMessages,
			DeletingMembers:  cmd.DeletingMembers,
			Muting:           cmd.Muting,
		})
		member = m
		return nil
	}); err != nil {
		return err
	}

	if err := h.emitter.Emit(events.MemberUpdatedEvent{
		ID:      member.ID(),
		GroupID: member.GroupID(),
		UserID:  member.UserID(),
		User: events.User{
			UserName:   member.User().Username(),
			HasPicture: member.User().HasPicture(),
		},
		DeletingMessages: member.DeletingMessages(),
		Muting:           member.Muting(),
		Adding:           member.Adding(),
		DeletingMembers:  member.DeletingMembers(),
		Admin:            member.Admin(),
	}); err != nil {
		return err
	}
	return nil
}
