package service_test

import (
	"errors"
	"fmt"
	"testing"

	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database/mock"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/service"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	mockqueue "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type InviteTestSuite struct {
	suite.Suite
	IDs     map[string]uuid.UUID
	Service service.ServiceLayer
}

func (s *InviteTestSuite) SetupSuite() {
	s.IDs = make(map[string]uuid.UUID)

	s.IDs["targetOK"] = uuid.New()
	s.IDs["targetNotFound"] = uuid.New()
	s.IDs["targetInGroup"] = uuid.New()
	s.IDs["targetInvited"] = uuid.New()

	s.IDs["issuerOK"] = uuid.New()
	s.IDs["issuerNotInGroup"] = uuid.New()
	s.IDs["issuerWithoutRights"] = uuid.New()

	s.IDs["groupOK"] = uuid.New()

	s.IDs["inviteOK"] = uuid.New()
	s.IDs["inviteOK2"] = uuid.New()
	s.IDs["inviteNotFound"] = uuid.New()
	s.IDs["inviteAnswered"] = uuid.New()

	db := new(mockdb.MockGroupsDB)
	db.On("GetUserByID", s.IDs["targetNotFound"]).Return(nil, errors.New("not found"))
	db.On("GetUserByID", s.IDs["targetInGroup"]).Return(&models.User{ID: s.IDs["targetInGroup"]}, nil)
	db.On("GetUserByID", s.IDs["targetInvited"]).Return(&models.User{ID: s.IDs["targetInvited"]}, nil)
	db.On("GetUserByID", s.IDs["targetOK"]).Return(&models.User{ID: s.IDs["targetOK"]}, nil)

	db.On("GetMemberByUserGroupID", s.IDs["issuerNotInGroup"], s.IDs["groupOK"]).Return(nil, errors.New("not found"))
	db.On("GetMemberByUserGroupID", s.IDs["issuerWithoutRights"], s.IDs["groupOK"]).Return(&models.Member{}, nil)
	db.On("GetMemberByUserGroupID", s.IDs["issuerOK"], s.IDs["groupOK"]).Return(&models.Member{Creator: true}, nil)
	db.On("GetMemberByUserGroupID", s.IDs["targetInGroup"], s.IDs["groupOK"]).Return(&models.Member{}, nil)
	db.On("GetMemberByUserGroupID", s.IDs["targetInvited"], s.IDs["groupOK"]).Return(nil, errors.New("not found"))
	db.On("GetMemberByUserGroupID", s.IDs["targetOK"], s.IDs["groupOK"]).Return(nil, errors.New("not found"))

	db.On("IsUserInvited", s.IDs["targetInvited"], s.IDs["groupOK"]).Return(true, nil)
	db.On("IsUserInvited", s.IDs["targetOK"], s.IDs["groupOK"]).Return(false, nil)

	db.On("CreateInvite", mock.Anything).Return(&models.Invite{ID: s.IDs["inviteOK"]}, nil)

	db.On("GetInviteByID", s.IDs["inviteNotFound"]).Return(nil, errors.New("not found"))
	db.On("GetInviteByID", s.IDs["inviteOK"]).Return(&models.Invite{TargetID: s.IDs["targetOK"], Status: models.INVITE_AWAITING}, nil)
	db.On("GetInviteByID", s.IDs["inviteOK2"]).Return(&models.Invite{TargetID: s.IDs["targetOK"], Status: models.INVITE_AWAITING}, nil)
	db.On("GetInviteByID", s.IDs["inviteAnswered"]).Return(&models.Invite{TargetID: s.IDs["targetOK"], Status: models.INVITE_DECLINE}, nil)

	db.On("UpdateInvite", mock.Anything).Return(&models.Invite{GroupID: s.IDs["groupOK"]}, nil)

	db.On("CreateMember", mock.Anything).Return(&models.Member{}, nil)

	db.On("GetGroupByID", s.IDs["groupOK"]).Return(&models.Group{}, nil)

	emitter := new(mockqueue.MockEmitter)
	emitter.On("Emit", mock.Anything).Return(nil)

	s.Service = service.NewService(db, nil, emitter)
}

func (s *InviteTestSuite) TestAddInvite() {
	testCases := []struct {
		desc           string
		issuerID       uuid.UUID
		targetID       uuid.UUID
		groupID        uuid.UUID
		expectedResult *models.Invite
		expectedError  error
	}{
		{
			desc:          "issuer not in group",
			issuerID:      s.IDs["issuerNotInGroup"],
			targetID:      uuid.New(),
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewNotFound("group not found"),
		},
		{
			desc:          "issuer has no rights",
			issuerID:      s.IDs["issuerWithoutRights"],
			targetID:      uuid.New(),
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewForbidden("user can't send invites to this group"),
		},
		{
			desc:          "target not found",
			issuerID:      s.IDs["issuerOK"],
			targetID:      s.IDs["targetNotFound"],
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewNotFound(fmt.Sprintf("user with ID %v not found", s.IDs["targetNotFound"])),
		},
		{
			desc:          "target in group",
			issuerID:      s.IDs["issuerOK"],
			targetID:      s.IDs["targetInGroup"],
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewForbidden("user is already a member of group"),
		},
		{
			desc:          "target invited",
			issuerID:      s.IDs["issuerOK"],
			targetID:      s.IDs["targetInvited"],
			groupID:       s.IDs["groupOK"],
			expectedError: apperrors.NewForbidden("user already invited"),
		},
		{
			desc:           "success",
			issuerID:       s.IDs["issuerOK"],
			targetID:       s.IDs["targetOK"],
			groupID:        s.IDs["groupOK"],
			expectedResult: &models.Invite{ID: s.IDs["inviteOK"]},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			res, err := s.Service.AddInvite(tC.issuerID, tC.targetID, tC.groupID)
			s.Equal(tC.expectedResult, res)
			s.Equal(tC.expectedError, err)
		})
	}
}

func (s *InviteTestSuite) TestRespondInvite() {
	testCases := []struct {
		desc           string
		userID         uuid.UUID
		inviteID       uuid.UUID
		response       bool
		expectedInvite *models.Invite
		expectedGroup  *models.Group
		expectedError  error
	}{
		{
			desc:          "invite not found",
			userID:        s.IDs["targetOK"],
			inviteID:      s.IDs["inviteNotFound"],
			expectedError: apperrors.NewNotFound("invite not found"),
		},
		{
			desc:          "invite not for target",
			userID:        s.IDs["issuerOK"],
			inviteID:      s.IDs["inviteOK"],
			expectedError: apperrors.NewNotFound("invite not found"),
		},
		{
			desc:          "invite answered",
			userID:        s.IDs["targetOK"],
			inviteID:      s.IDs["inviteAnswered"],
			expectedError: apperrors.NewConflict("invite already answered"),
		},
		{
			desc:           "invite declined",
			userID:         s.IDs["targetOK"],
			inviteID:       s.IDs["inviteOK"],
			response:       false,
			expectedInvite: &models.Invite{GroupID: s.IDs["groupOK"]},
		},
		{
			desc:     "invite accepted",
			userID:   s.IDs["targetOK"],
			inviteID: s.IDs["inviteOK2"], // second inviteOK is needed because there invite value returned by mock is a pointer which is updated in
			// another call - "invite declined" therefore if we would use the same invite ID it method GetInviteByID would return an object that was
			// overwritten
			response:       true,
			expectedInvite: &models.Invite{GroupID: s.IDs["groupOK"]},
			expectedGroup:  &models.Group{},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			invite, group, err := s.Service.RespondInvite(tC.userID, tC.inviteID, tC.response)
			s.Equal(tC.expectedInvite, invite)
			s.Equal(tC.expectedGroup, group)
			s.Equal(tC.expectedError, err)
		})
	}
}

func TestInviteSuite(t *testing.T) {
	suite.Run(t, &InviteTestSuite{})
}
