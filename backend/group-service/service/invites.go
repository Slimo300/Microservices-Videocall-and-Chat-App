package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/google/uuid"
)

func (srv *GroupService) AddInvite(ctx context.Context, userID, targetUserID, groupID uuid.UUID) (*models.Invite, error) {

	member, err := srv.DB.GetMemberByUserGroupID(ctx, userID, groupID)
	if err != nil {
		return nil, apperrors.NewNotFound("group not found")
	}

	if !member.Adding && !member.Admin && !member.Creator {
		return nil, apperrors.NewForbidden("user can't send invites to this group")
	}

	targetUser, err := srv.DB.GetUserByID(ctx, targetUserID)
	if err != nil {
		return nil, apperrors.NewNotFound(fmt.Sprintf("user with ID %v not found", targetUserID))
	}

	_, err = srv.DB.GetMemberByUserGroupID(ctx, targetUserID, groupID)
	if err == nil {
		return nil, apperrors.NewForbidden("user is already a member of group")
	}

	isInvited, err := srv.DB.IsUserInvited(ctx, targetUserID, groupID)
	if err != nil || isInvited {
		return nil, apperrors.NewForbidden("user already invited")
	}

	invite, err := srv.DB.CreateInvite(ctx, &models.Invite{
		ID:       uuid.New(),
		IssId:    userID,
		TargetID: targetUser.ID,
		GroupID:  groupID,
		Status:   models.INVITE_AWAITING,
		Created:  time.Now(),
		Modified: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if err := srv.Emitter.Emit(events.InviteSentEvent{
		ID:       invite.ID,
		IssuerID: invite.IssId,
		TargetID: invite.TargetID,
		GroupID:  invite.GroupID,
	}); err != nil {
		return nil, err
	}

	return invite, nil
}

func (srv *GroupService) RespondInvite(ctx context.Context, userID, inviteID uuid.UUID, answer bool) (*models.Invite, *models.Group, error) {

	invite, err := srv.DB.GetInviteByID(ctx, inviteID)
	if err != nil {
		return nil, nil, apperrors.NewNotFound("invite not found")
	}

	if invite.TargetID != userID {
		return nil, nil, apperrors.NewNotFound("invite not found")
	}

	if invite.Status != models.INVITE_AWAITING {
		return nil, nil, apperrors.NewConflict("invite already answered")
	}

	invite.Modified = time.Now()

	if !answer {
		invite.Status = models.INVITE_DECLINE
		invite, err = srv.DB.UpdateInvite(ctx, invite)
		if err != nil {
			return nil, nil, err
		}
		return invite, nil, nil
	}

	invite.Status = models.INVITE_ACCEPT
	invite, err = srv.DB.UpdateInvite(ctx, invite)
	if err != nil {
		return nil, nil, err
	}

	member, err := srv.DB.CreateMember(ctx, &models.Member{
		ID:      uuid.New(),
		GroupID: invite.GroupID,
		UserID:  userID,
	})
	if err != nil {
		return nil, nil, err
	}

	group, err := srv.DB.GetGroupByID(ctx, invite.GroupID)
	if err != nil {
		return nil, nil, err
	}

	if err = srv.Emitter.Emit(events.MemberCreatedEvent{
		ID:      member.ID,
		GroupID: member.GroupID,
		UserID:  member.UserID,
		User: events.User{
			UserName:   member.User.UserName,
			HasPicture: member.User.HasPicture,
		},
		Creator: member.Creator,
	}); err != nil {
		return nil, nil, err
	}

	if err = srv.Emitter.Emit(events.InviteRespondedEvent{
		ID:       invite.ID,
		IssuerID: invite.IssId,
		TargetID: invite.TargetID,
		Target: events.User{
			UserName:   invite.Target.UserName,
			HasPicture: invite.Target.HasPicture,
		},
		GroupID: invite.GroupID,
		Group: events.Group{
			Name:       invite.Group.Name,
			HasPicture: invite.Group.HasPicture,
		},
		Status: int(invite.Status),
	}); err != nil {
		return nil, nil, err
	}

	return invite, group, nil
}

func (srv *GroupService) GetUserInvites(ctx context.Context, userID uuid.UUID, num, offset int) ([]*models.Invite, error) {
	return srv.DB.GetUserInvites(ctx, userID, num, offset)
}
