package service

import (
	"context"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/google/uuid"
)

func (srv GroupService) CreateGroup(ctx context.Context, userID uuid.UUID, name string) (*models.Group, error) {

	group, err := srv.DB.CreateGroup(ctx, &models.Group{
		ID:      uuid.New(),
		Name:    name,
		Created: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	member, err := srv.DB.CreateMember(ctx, &models.Member{
		ID:      uuid.New(),
		GroupID: group.ID,
		UserID:  userID,
		Creator: true,
	})
	if err != nil {
		return nil, err
	}

	group.Members = append(group.Members, *member)

	if err := srv.Emitter.Emit(events.MemberCreatedEvent{
		ID:      member.ID,
		GroupID: member.GroupID,
		UserID:  member.UserID,
		Creator: member.Creator,
		User: events.User{
			UserName:   member.User.UserName,
			HasPicture: member.User.HasPicture,
		},
	}); err != nil {
		return nil, err
	}

	return group, nil
}

func (srv GroupService) DeleteGroup(ctx context.Context, userID, groupID uuid.UUID) (*models.Group, error) {
	member, err := srv.DB.GetMemberByUserGroupID(ctx, userID, groupID)
	if err != nil {
		return nil, apperrors.NewNotFound("group not found")
	}

	if !member.Creator {
		return nil, apperrors.NewForbidden("member can't delete this group")
	}

	group, err := srv.DB.DeleteGroup(ctx, groupID)
	if err != nil {
		return nil, apperrors.NewNotFound("group not found")
	}

	return group, nil
}

func (srv GroupService) GetUserGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error) {

	groups, err := srv.DB.GetUserGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
