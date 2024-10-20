package command

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/google/uuid"
)

type CreateGroupCommand struct {
	UserID uuid.UUID
	Name   string
}

type CreateGroupHandler struct {
	repo    database.GroupsRepository
	emitter msgqueue.EventEmiter
}

func NewCreateGroupHandler(repo database.GroupsRepository, emitter msgqueue.EventEmiter) CreateGroupHandler {
	if repo == nil {
		panic("repo is nil")
	}
	if emitter == nil {
		panic("emitter is nil")
	}
	return CreateGroupHandler{repo: repo, emitter: emitter}
}

func (h CreateGroupHandler) Handle(ctx context.Context, cmd CreateGroupCommand) error {
	group := models.CreateGroup(cmd.UserID, cmd.Name)
	group, err := h.repo.CreateGroup(ctx, group)
	if err != nil {
		return err
	}

	member := group.Members()[0]
	if err := h.emitter.Emit(events.MemberCreatedEvent{
		ID:      member.ID(),
		GroupID: member.GroupID(),
		UserID:  member.UserID(),
		Creator: member.Creator(),
		User: events.User{
			UserName:   member.User().Username(),
			HasPicture: member.User().HasPicture(),
		},
	}); err != nil {
		return err
	}

	return nil
}
