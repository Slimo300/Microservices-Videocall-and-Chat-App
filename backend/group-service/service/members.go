package service

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/google/uuid"
)

func (srv *GroupService) GrantRights(ctx context.Context, userID, memberID uuid.UUID, rights models.MemberRights) (*models.Member, error) {
	member, err := srv.DB.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, apperrors.NewNotFound("member not found")
	}

	group, err := srv.DB.GetGroupByID(ctx, member.GroupID)
	if err != nil {
		return nil, apperrors.NewNotFound("member not found")
	}

	issuerMember, err := srv.DB.GetMemberByUserGroupID(ctx, userID, group.ID)
	if err != nil {
		return nil, apperrors.NewNotFound("member not found")
	}

	if !issuerMember.CanAlter(*member) {
		return nil, apperrors.NewForbidden("user cannot set rights")
	}

	member.ApplyRights(rights)
	member, err = srv.DB.UpdateMember(ctx, member)
	if err != nil {
		return nil, err
	}

	if err := srv.Emitter.Emit(events.MemberUpdatedEvent{
		ID:      member.ID,
		GroupID: member.GroupID,
		UserID:  member.UserID,
		User: events.User{
			UserName:   member.User.UserName,
			HasPicture: member.User.HasPicture,
		},
		DeletingMessages: member.DeletingMessages,
		Muting:           member.Muting,
		Adding:           member.Adding,
		DeletingMembers:  member.DeletingMembers,
		Admin:            member.Admin,
	}); err != nil {
		return nil, err
	}

	return member, nil
}

func (srv *GroupService) DeleteMember(ctx context.Context, userID, memberID uuid.UUID) (*models.Member, error) {
	member, err := srv.DB.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, apperrors.NewNotFound("member not found")
	}

	group, err := srv.DB.GetGroupByID(ctx, member.GroupID)
	if err != nil {
		return nil, apperrors.NewNotFound("member not found")
	}

	issuerMember, err := srv.DB.GetMemberByUserGroupID(ctx, userID, group.ID)
	if err != nil {
		return nil, apperrors.NewNotFound("member not found")
	}

	if !issuerMember.CanDelete(*member) {
		return nil, apperrors.NewForbidden("user can't delete from this group")
	}

	member, err = srv.DB.DeleteMember(ctx, member.ID)
	if err != nil {
		return nil, err
	}

	if err := srv.Emitter.Emit(events.MemberDeletedEvent{
		ID:      member.ID,
		UserID:  member.UserID,
		GroupID: member.GroupID,
	}); err != nil {
		return nil, err
	}

	return member, nil
}
