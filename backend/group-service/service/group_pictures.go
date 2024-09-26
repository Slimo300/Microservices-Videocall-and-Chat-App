package service

import (
	"context"
	"mime/multipart"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/google/uuid"
)

func (srv *GroupService) SetGroupPicture(ctx context.Context, userID, groupID uuid.UUID, file multipart.File) error {

	group, err := srv.DB.GetGroupByID(ctx, groupID)
	if err != nil {
		return apperrors.NewNotFound("group not found")
	}
	member, err := srv.DB.GetMemberByUserGroupID(ctx, userID, groupID)
	if err != nil || !member.Creator {
		return apperrors.NewForbidden("user has no rights to set group picture")
	}
	if !group.HasPicture {
		group.HasPicture = true
		_, err = srv.DB.UpdateGroup(ctx, group)
		if err != nil {
			return err
		}
	}
	if err := srv.Storage.UploadFile(ctx, groupID.String(), file, storage.PUBLIC_READ); err != nil {
		return err
	}
	return nil
}

func (srv *GroupService) DeleteGroupPicture(ctx context.Context, userID, groupID uuid.UUID) error {
	group, err := srv.DB.GetGroupByID(ctx, groupID)
	if err != nil {
		return apperrors.NewNotFound("group not found")
	}
	if !group.HasPicture {
		return apperrors.NewBadRequest("group has no picture")
	}

	membership, err := srv.DB.GetMemberByUserGroupID(ctx, userID, groupID)
	if err != nil || !membership.Creator {
		return apperrors.NewForbidden("user has no rights to delete group picture")
	}

	group.HasPicture = false
	if _, err := srv.DB.UpdateGroup(ctx, group); err != nil {
		return err
	}

	if err := srv.Storage.DeleteFile(ctx, group.ID.String()); err != nil {
		return err
	}

	return nil
}
