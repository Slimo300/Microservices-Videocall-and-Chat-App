package service_test

import (
	"errors"
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

type MembersTestSuite struct {
	suite.Suite
	IDs     map[string]uuid.UUID
	Service service.ServiceLayer
}

func (s *MembersTestSuite) SetupSuite() {
	s.IDs = make(map[string]uuid.UUID)

	s.IDs["groupOK"] = uuid.New()
	s.IDs["groupNotFound"] = uuid.New()

	s.IDs["userOK"] = uuid.New()
	s.IDs["userOutsideGroup"] = uuid.New()
	s.IDs["userWithoutRights"] = uuid.New()

	s.IDs["memberOK"] = uuid.New()
	s.IDs["memberWithoutGroup"] = uuid.New()
	s.IDs["memberNotFound"] = uuid.New()

	db := new(mockdb.MockGroupsDB)
	db.On("GetMemberByID", s.IDs["memberOK"]).Return(&models.Member{GroupID: s.IDs["groupOK"]}, nil)
	db.On("GetMemberByID", s.IDs["memberNotFound"]).Return(nil, errors.New("not found"))
	db.On("GetMemberByID", s.IDs["memberWithoutGroup"]).Return(&models.Member{GroupID: s.IDs["groupNotFound"]}, nil)

	db.On("GetGroupByID", s.IDs["groupOK"]).Return(&models.Group{ID: s.IDs["groupOK"]}, nil)
	db.On("GetGroupByID", s.IDs["groupNotFound"]).Return(nil, errors.New("not found"))

	db.On("GetMemberByUserGroupID", s.IDs["userOutsideGroup"], s.IDs["groupOK"]).Return(nil, errors.New("not found"))
	db.On("GetMemberByUserGroupID", s.IDs["userWithoutRights"], s.IDs["groupOK"]).Return(&models.Member{ID: uuid.New()}, nil) // uuid.New() is necessary, otherwise members would have the same ID and the program would allow to delete since users can delete themselves
	db.On("GetMemberByUserGroupID", s.IDs["userOK"], s.IDs["groupOK"]).Return(&models.Member{ID: s.IDs["memberOK"], Admin: true}, nil)

	db.On("DeleteMember", mock.Anything).Return(&models.Member{ID: s.IDs["memberOK"]}, nil)

	db.On("UpdateMember", mock.Anything).Return(&models.Member{ID: s.IDs["memberOK"]}, nil)

	emitter := new(mockqueue.MockEmitter)
	emitter.On("Emit", mock.Anything).Return(nil)

	s.Service = service.NewService(db, nil, emitter)
}

func (s *MembersTestSuite) TestDeleteMember() {
	testCases := []struct {
		desc           string
		userID         uuid.UUID
		memberID       uuid.UUID
		expectedResult *models.Member
		expectedError  error
	}{
		{
			desc:          "member not found",
			userID:        s.IDs["userOK"],
			memberID:      s.IDs["memberNotFound"],
			expectedError: apperrors.NewNotFound("member not found"),
		},
		{
			desc:          "member not found",
			userID:        s.IDs["userOK"],
			memberID:      s.IDs["memberWithoutGroup"],
			expectedError: apperrors.NewNotFound("member not found"),
		},
		{
			desc:          "issuer not in group",
			userID:        s.IDs["userOutsideGroup"],
			memberID:      s.IDs["memberOK"],
			expectedError: apperrors.NewNotFound("member not found"),
		},
		{
			desc:          "issuer has no rights",
			userID:        s.IDs["userWithoutRights"],
			memberID:      s.IDs["memberOK"],
			expectedError: apperrors.NewForbidden("user can't delete from this group"),
		},
		{
			desc:          "issuer has no rights",
			userID:        s.IDs["userWithoutRights"],
			memberID:      s.IDs["memberOK"],
			expectedError: apperrors.NewForbidden("user can't delete from this group"),
		},
		{
			desc:           "success",
			userID:         s.IDs["userOK"],
			memberID:       s.IDs["memberOK"],
			expectedResult: &models.Member{ID: s.IDs["memberOK"]},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			res, err := s.Service.DeleteMember(tC.userID, tC.memberID)
			s.Equal(tC.expectedResult, res)
			s.Equal(tC.expectedError, err)
		})
	}
}

func (s *MembersTestSuite) TestGrantRights() {
	testCases := []struct {
		desc           string
		userID         uuid.UUID
		memberID       uuid.UUID
		expectedResult *models.Member
		expectedError  error
	}{

		{
			desc:          "member not found",
			userID:        s.IDs["userOK"],
			memberID:      s.IDs["memberNotFound"],
			expectedError: apperrors.NewNotFound("member not found"),
		},
		{
			desc:          "member not found",
			userID:        s.IDs["userOK"],
			memberID:      s.IDs["memberWithoutGroup"],
			expectedError: apperrors.NewNotFound("member not found"),
		},
		{
			desc:          "issuer not in group",
			userID:        s.IDs["userOutsideGroup"],
			memberID:      s.IDs["memberOK"],
			expectedError: apperrors.NewNotFound("member not found"),
		},
		{
			desc:          "issuer has no rights",
			userID:        s.IDs["userWithoutRights"],
			memberID:      s.IDs["memberOK"],
			expectedError: apperrors.NewForbidden("user cannot set rights"),
		},
		{
			desc:           "success",
			userID:         s.IDs["userOK"],
			memberID:       s.IDs["memberOK"],
			expectedResult: &models.Member{ID: s.IDs["memberOK"]},
		},
		{
			desc:           "success",
			userID:         s.IDs["userOK"],
			memberID:       s.IDs["memberOK"],
			expectedResult: &models.Member{ID: s.IDs["memberOK"]},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			res, err := s.Service.GrantRights(tC.userID, tC.memberID, models.MemberRights{})
			s.Equal(tC.expectedResult, res)
			s.Equal(tC.expectedError, err)
		})
	}
}

func TestMembersSuite(t *testing.T) {
	suite.Run(t, &MembersTestSuite{})
}
