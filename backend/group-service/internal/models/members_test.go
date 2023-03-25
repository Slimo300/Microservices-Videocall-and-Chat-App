package models_test

import (
	"testing"

	"github.com/Slimo300/chat-groupservice/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MemberTestSuite struct {
	suite.Suite
	basic    models.Member
	basic2   models.Member
	basic3   models.Member
	deleter  models.Member
	deleter2 models.Member
	deleter3 models.Member
	admin    models.Member
	admin2   models.Member
	admin3   models.Member
	creator  models.Member
	creator2 models.Member
}

func (s *MemberTestSuite) SetupSuite() {
	s.basic = models.Member{ID: uuid.New()}
	s.basic2 = models.Member{ID: uuid.New()}
	s.basic3 = models.Member{ID: uuid.New()}
	s.deleter = models.Member{ID: uuid.New(), DeletingMembers: true}
	s.deleter3 = models.Member{ID: uuid.New(), DeletingMembers: true}
	s.admin = models.Member{ID: uuid.New(), Admin: true}
	s.admin2 = models.Member{ID: uuid.New(), Admin: true}
	s.admin3 = models.Member{ID: uuid.New(), Admin: true}
	s.creator = models.Member{ID: uuid.New(), Creator: true}
	s.creator2 = models.Member{ID: uuid.New(), Creator: true}
}

func (s *MemberTestSuite) TestCanDelete() {
	s.True(s.creator.CanDelete(s.basic))
	s.True(s.creator.CanDelete(s.deleter))
	s.True(s.creator.CanDelete(s.admin))
	s.False(s.creator.CanDelete(s.creator2))

	s.True(s.admin.CanDelete(s.basic))
	s.True(s.admin.CanDelete(s.deleter))
	s.False(s.admin.CanDelete(s.admin2))
	s.False(s.admin.CanDelete(s.creator))

	s.True(s.deleter.CanDelete(s.basic))
	s.True(s.deleter.CanDelete(s.deleter2))
	s.False(s.deleter.CanDelete(s.admin))
	s.False(s.deleter.CanDelete(s.creator))

	s.False(s.basic.CanDelete(s.basic2))
	s.False(s.basic.CanDelete(s.deleter))
	s.False(s.basic.CanDelete(s.admin))
	s.False(s.basic.CanDelete(s.creator))

	s.True(s.basic.CanDelete(s.basic))
	s.True(s.deleter.CanDelete(s.deleter))
	s.True(s.admin.CanDelete(s.admin))
	s.False(s.creator.CanDelete(s.creator))
}

func (s *MemberTestSuite) TestCanAlter() {
	s.True(s.creator.CanAlter(s.basic))
	s.True(s.creator.CanAlter(s.deleter))
	s.True(s.creator.CanAlter(s.admin))
	s.False(s.creator.CanAlter(s.creator))

	s.True(s.admin.CanAlter(s.basic))
	s.True(s.admin.CanAlter(s.deleter))
	s.False(s.admin.CanAlter(s.admin))
	s.False(s.admin.CanAlter(s.creator))

	s.False(s.deleter.CanAlter(s.basic))
	s.False(s.deleter.CanAlter(s.deleter))
	s.False(s.deleter.CanAlter(s.admin))
	s.False(s.deleter.CanAlter(s.creator))

	s.False(s.basic.CanAlter(s.basic))
	s.False(s.basic.CanAlter(s.deleter))
	s.False(s.basic.CanAlter(s.admin))
	s.False(s.basic.CanAlter(s.creator))
}

func (s *MemberTestSuite) TestApplyRights() {
	s.False(s.basic3.Adding)

	if err := s.basic3.ApplyRights(models.MemberRights{
		Adding: 1,
	}); err != nil {
		s.Fail(err.Error())
	}
	s.True(s.basic3.Adding)

	if err := s.basic3.ApplyRights(models.MemberRights{
		Adding: -1,
	}); err != nil {
		s.Fail(err.Error())
	}
	s.False(s.basic3.Adding)

	if err := s.basic3.ApplyRights(models.MemberRights{
		Adding: 0,
	}); err != nil {
		s.Fail(err.Error())
	}
	s.False(s.basic3.Adding)

	if err := s.admin3.ApplyRights(models.MemberRights{
		Adding:           1,
		DeletingMessages: 0,
		DeletingMembers:  -1,
		Admin:            -1,
	}); err != nil {
		s.Fail(err.Error())
	}
	s.True(s.admin3.Adding)
	s.False(s.admin3.DeletingMessages)
	s.False(s.admin3.DeletingMembers)
	s.False(s.admin3.Admin)

	s.Error(s.admin3.ApplyRights(models.MemberRights{Adding: 2}))

}

func TestMembers(t *testing.T) {
	suite.Run(t, &MemberTestSuite{})
}
