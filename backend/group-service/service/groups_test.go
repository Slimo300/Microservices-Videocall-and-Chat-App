package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database/mock"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	mockqueue "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/mock"
)

type GroupsTestSuite struct {
	suite.Suite
	IDs     map[string]uuid.UUID
	Service service.Service
}

func (s *GroupsTestSuite) SetupSuite() {

	s.IDs = make(map[string]uuid.UUID)

	s.IDs["userOK"] = uuid.New()
	s.IDs["groupOK"] = uuid.New()
	s.IDs["memberOK"] = uuid.New()
	s.IDs["userNotInGroup"] = uuid.New()
	s.IDs["userWithoutRights"] = uuid.New()

	db := new(mockdb.GroupsMockRepository)
	db.On("GetMemberByUserGroupID", mock.Anything, s.IDs["userOK"], s.IDs["groupOK"]).Return(&models.Member{ID: uuid.New(), Creator: true}, nil)
	db.On("GetMemberByUserGroupID", mock.Anything, s.IDs["userNotInGroup"], s.IDs["groupOK"]).Return(nil, errors.New("group not found"))
	db.On("GetMemberByUserGroupID", mock.Anything, s.IDs["userWithoutRights"], s.IDs["groupOK"]).Return(&models.Member{ID: uuid.New(), Creator: false}, nil)
	db.On("DeleteGroup", mock.Anything, s.IDs["groupOK"]).Return(&models.Group{ID: s.IDs["groupOK"]}, nil)

	db.On("CreateGroup", mock.Anything, mock.AnythingOfType("*models.Group")).Return(&models.Group{ID: s.IDs["groupOK"]}, nil)
	db.On("CreateMember", mock.Anything, mock.AnythingOfType("*models.Member")).Return(&models.Member{ID: s.IDs["memberOK"]}, nil)

	emitter := new(mockqueue.MockEmitter)
	emitter.On("Emit", mock.Anything).Return(nil)

	s.Service = service.NewService(db, nil, emitter)
}

func (s *GroupsTestSuite) TestDeleteGroup() {
	testCases := []struct {
		desc           string
		userID         uuid.UUID
		groupID        uuid.UUID
		expectedResult *models.Group
		expectedError  error
	}{
		{
			desc:          "member not found",
			userID:        s.IDs["userNotInGroup"],
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewNotFound("group not found"),
		},
		{
			desc:          "member is not a creator",
			userID:        s.IDs["userWithoutRights"],
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewForbidden("member can't delete this group"),
		},
		{
			desc:           "success",
			userID:         s.IDs["userOK"],
			groupID:        s.IDs["groupOK"],
			expectedResult: &models.Group{ID: s.IDs["groupOK"]},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			group, err := s.Service.DeleteGroup(context.Background(), tC.userID, tC.groupID)
			s.Equal(tC.expectedError, err)
			s.Equal(tC.expectedResult, group)
		})
	}
}

func (s *GroupsTestSuite) TestCreateGroup() {
	group, err := s.Service.CreateGroup(context.Background(), s.IDs["userOK"], "New Group")
	s.Equal(nil, err)
	s.Equal(s.IDs["groupOK"], group.ID)
	s.Equal(s.IDs["memberOK"], group.Members[0].ID)
}

func TestGroupsSuite(t *testing.T) {
	suite.Run(t, &GroupsTestSuite{})
}
