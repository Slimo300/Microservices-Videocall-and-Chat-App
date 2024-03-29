package service

import (
	"mime/multipart"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/google/uuid"
)

func (srv *GroupService) SetGroupPicture(userID, groupID uuid.UUID, file multipart.File) (string, error) {

	group, err := srv.DB.GetGroupByID(groupID)
	if err != nil {
		return "", apperrors.NewNotFound("group not found")
	}

	membership, err := srv.DB.GetMemberByUserGroupID(userID, groupID)
	if err != nil || !membership.Creator {
		return "", apperrors.NewForbidden("user has no rights to set group picture")
	}

	if group.Picture == "" {
		group.Picture = uuid.NewString()
		group, err = srv.DB.UpdateGroup(group)
		if err != nil {
			return "", err
		}
	}

	if err := srv.Storage.UploadFile(file, group.Picture); err != nil {
		return "", err
	}

	return group.Picture, nil
}

func (srv *GroupService) DeleteGroupPicture(userID, groupID uuid.UUID) error {
	group, err := srv.DB.GetGroupByID(groupID)
	if err != nil {
		return apperrors.NewNotFound("group not found")
	}

	membership, err := srv.DB.GetMemberByUserGroupID(userID, groupID)
	if err != nil || !membership.Creator {
		return apperrors.NewForbidden("user has no rights to delete group picture")
	}

	if group.Picture == "" {
		return apperrors.NewBadRequest("group has no picture")
	}

	if err := srv.Storage.DeleteFile(group.Picture); err != nil {
		return err
	}

	return nil
}
