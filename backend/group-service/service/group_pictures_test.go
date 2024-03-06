package service_test

import (
	"errors"
	"testing"

	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database/mock"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/service"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"

	"github.com/google/uuid"
)

type GroupPicturesTestSuite struct {
	suite.Suite
	IDs     map[string]uuid.UUID
	Service service.ServiceLayer
}

func (s *GroupPicturesTestSuite) SetupSuite() {

	s.IDs = make(map[string]uuid.UUID)

	s.IDs["userOK"] = uuid.New()
	s.IDs["groupOK"] = uuid.New()
	s.IDs["userWithoutRights"] = uuid.New()
	s.IDs["groupWithoutPicture"] = uuid.New()

	db := new(mockdb.MockGroupsDB)
	db.On("GetMemberByUserGroupID", s.IDs["userOK"], s.IDs["groupOK"]).Return(&models.Member{ID: uuid.New(), Creator: true}, nil)
	db.On("GetMemberByUserGroupID", s.IDs["userOK"], s.IDs["groupWithoutPicture"]).Return(&models.Member{ID: uuid.New(), Creator: true}, nil)
	db.On("GetMemberByUserGroupID", s.IDs["userWithoutRights"], s.IDs["groupOK"]).Return(&models.Member{ID: uuid.New(), Creator: false}, nil)
	db.On("GetMemberByUserGroupID", s.IDs["userOK"], mock.Anything).Return(nil, errors.New("member not found"))

	db.On("GetGroupByID", s.IDs["groupOK"]).Return(&models.Group{ID: s.IDs["groupOK"], Picture: "picture"}, nil)
	db.On("GetGroupByID", s.IDs["groupWithoutPicture"]).Return(&models.Group{ID: s.IDs["groupWithoutPicture"]}, nil)
	db.On("GetGroupByID", mock.Anything).Return(nil, errors.New("group not found"))

	db.On("UpdateGroup", mock.Anything).Return(&models.Group{ID: s.IDs["groupOK"], Picture: "picture"}, nil)

	storage := new(storage.MockStorage)

	storage.On("DeleteFile", "force_error").Return(errors.New("storage error"))
	storage.On("DeleteFile", mock.Anything).Return(nil)
	storage.On("UploadFile", mock.Anything, mock.Anything).Return(nil)

	s.Service = service.NewService(db, storage, nil)
}

func (s *GroupPicturesTestSuite) TestDeleteGroupPicture() {

	testCases := []struct {
		desc          string
		userID        uuid.UUID
		groupID       uuid.UUID
		expectedError error
	}{
		{
			desc:          "group not found",
			userID:        s.IDs["userOK"],
			groupID:       uuid.New(),
			expectedError: apperrors.NewNotFound("group not found"),
		},
		{
			desc:          "no rigts to set",
			userID:        s.IDs["userWithoutRights"],
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewForbidden("user has no rights to delete group picture"),
		},
		{
			desc:          "group has no picture",
			userID:        s.IDs["userOK"],
			groupID:       s.IDs["groupWithoutPicture"],
			expectedError: apperrors.NewBadRequest("group has no picture"),
		},
		{
			desc:          "group has no picture",
			userID:        s.IDs["userOK"],
			groupID:       s.IDs["groupWithoutPicture"],
			expectedError: apperrors.NewBadRequest("group has no picture"),
		},
		{
			desc:          "success",
			userID:        s.IDs["userOK"],
			groupID:       s.IDs["groupOK"],
			expectedError: nil,
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.Service.DeleteGroupPicture(tC.userID, tC.groupID)
			s.Equal(tC.expectedError, err)
		})
	}
}

func (s *GroupPicturesTestSuite) TestSetGroupPicture() {
	testCases := []struct {
		desc           string
		userID         uuid.UUID
		groupID        uuid.UUID
		expectedResult string
		expectedError  error
	}{
		{
			desc:          "group not found",
			userID:        s.IDs["userOK"],
			groupID:       uuid.New(),
			expectedError: apperrors.NewNotFound("group not found"),
		},
		{
			desc:          "no rigts to set",
			userID:        s.IDs["userWithoutRights"],
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewForbidden("user has no rights to set group picture"),
		},
		{
			desc:           "generate picture ID",
			userID:         s.IDs["userOK"],
			groupID:        s.IDs["groupWithoutPicture"],
			expectedResult: "picture",
			expectedError:  nil,
		},
		{
			desc:           "success",
			userID:         s.IDs["userOK"],
			groupID:        s.IDs["groupOK"],
			expectedResult: "picture",
			expectedError:  nil,
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			pictureID, err := s.Service.SetGroupPicture(tC.userID, tC.groupID, nil)
			s.Equal(tC.expectedResult, pictureID)
			s.Equal(tC.expectedError, err)
		})
	}
}

func TestGroupPicturesSuite(t *testing.T) {
	suite.Run(t, &GroupPicturesTestSuite{})
}
