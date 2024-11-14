package database_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database/orm"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// func TestMain(m *testing.M) {

// 	os.Exit(m.Run())
// }

type MessagesSuite struct {
	suite.Suite
	repo database.MessagesRepository
}

func TestSuites(t *testing.T) {
	gormRepository, err := orm.NewMessagesGormRepository("root:secret@(localhost:3306)/app")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []*MessagesSuite{{repo: gormRepository}} {
		s := s
		suite.Run(t, s)
	}
}

func (s *MessagesSuite) TestMembers() {
	s.T().Parallel()
	ctx := context.Background()
	member := models.NewMember(uuid.New(), uuid.New(), uuid.New(), uuid.NewString(), false)

	_, err := s.repo.GetUserGroupMember(ctx, member.UserID(), member.GroupID())
	s.Require().Error(err)
	s.Require().NoError(s.repo.CreateMember(ctx, member))
	returnedMember, err := s.repo.GetUserGroupMember(ctx, member.UserID(), member.GroupID())
	s.Require().NoError(err)

	s.Equal(returnedMember.DeletingMessages(), false)
	s.Require().NoError(s.repo.UpdateMember(ctx, member.ID(), func(m *models.Member) bool {
		return m.UpdateRights(false, true)
	}))
	returnedMember, err = s.repo.GetUserGroupMember(ctx, member.UserID(), member.GroupID())
	s.Require().NoError(err)
	s.Equal(returnedMember.DeletingMessages(), true)

	s.Require().NoError(s.repo.DeleteMember(ctx, member.ID()))
	_, err = s.repo.GetUserGroupMember(ctx, member.UserID(), member.GroupID())
	s.Error(err)
}

func (s *MessagesSuite) TestGetGroupMessages() {
	s.T().Parallel()
	ctx := context.Background()
	member := models.NewMember(uuid.New(), uuid.New(), uuid.New(), uuid.NewString(), false)

	message1 := models.NewMessage(uuid.New(), member.GroupID(), member.ID(), uuid.NewString(), time.Now(), nil)
	message2 := models.NewMessage(uuid.New(), member.GroupID(), member.ID(), uuid.NewString(), time.Now(), nil)
	message3 := models.NewMessage(uuid.New(), member.GroupID(), member.ID(), uuid.NewString(), time.Now(), nil)
	messageFile1 := models.NewMessageFile(message3.ID(), uuid.NewString(), "png")
	messageFile2 := models.NewMessageFile(message3.ID(), uuid.NewString(), "png")
	message3.AddFiles(messageFile1, messageFile2)

	s.Require().NoError(s.repo.CreateMember(ctx, member))
	s.Require().NoError(s.repo.CreateMessage(ctx, message1))
	s.Require().NoError(s.repo.CreateMessage(ctx, message2))
	s.Require().NoError(s.repo.CreateMessage(ctx, message3))

	msgs, err := s.repo.GetGroupMessages(ctx, member.UserID(), member.GroupID(), 0, 4)
	s.Require().NoError(err)
	s.Len(msgs, 3)
	s.True(findMessages(s.T(), msgs, message1.ID(), message2.ID(), message3.ID()))
}

func findMessages(t *testing.T, messages []models.Message, messageIDs ...uuid.UUID) bool {
	t.Helper()
	if messages == nil {
		return false
	}
	searchedMessages := make(map[uuid.UUID]bool)
	for i := range messageIDs {
		searchedMessages[messageIDs[i]] = true
	}
	for i := range messages {
		delete(searchedMessages, messages[i].ID())
	}
	return len(searchedMessages) == 0
}

func (s *MessagesSuite) TestDeleteMessageForYourself() {
	s.T().Parallel()
	ctx := context.Background()

	member1 := models.NewMember(uuid.New(), uuid.New(), uuid.New(), uuid.NewString(), true)
	member2 := models.NewMember(uuid.New(), uuid.New(), uuid.New(), uuid.NewString(), true)
	message1 := models.NewMessage(uuid.New(), member1.GroupID(), member1.ID(), uuid.NewString(), time.Now(), nil)
	message1.AddDeleters(member1)
	message2 := models.NewMessage(uuid.New(), member2.GroupID(), member2.ID(), uuid.NewString(), time.Now(), nil)
	s.Require().NoError(s.repo.CreateMember(ctx, member1))
	s.Require().NoError(s.repo.CreateMember(ctx, member2))
	s.Require().NoError(s.repo.CreateMessage(ctx, message1))
	s.Require().NoError(s.repo.CreateMessage(ctx, message2))

	testCases := []struct {
		desc           string
		userID         uuid.UUID
		messageID      uuid.UUID
		expectedResult error
	}{
		{
			desc:           "message_not_found",
			userID:         member1.UserID(),
			messageID:      member1.ID(),
			expectedResult: apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", member1.ID().String())),
		},
		{
			desc:           "user_not_in_group",
			userID:         member2.UserID(),
			messageID:      message1.ID(),
			expectedResult: apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", message1.ID().String())),
		},
		{
			desc:           "user_is_already_a_deleter",
			userID:         member1.UserID(),
			messageID:      message1.ID(),
			expectedResult: apperrors.NewConflict(fmt.Sprintf("message %s already deleted", message1.ID())),
		},
		{
			desc:           "success",
			userID:         member2.UserID(),
			messageID:      message2.ID(),
			expectedResult: nil,
		},
	}

	for _, tC := range testCases {
		s.T().Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.expectedResult, s.repo.DeleteMessageForYourself(ctx, tC.userID, tC.messageID))
		})
	}
}

func (s *MessagesSuite) TestDeleteMessageForEveryone() {
	s.T().Parallel()
	ctx := context.Background()
	member1 := models.NewMember(uuid.New(), uuid.New(), uuid.New(), uuid.NewString(), true)
	member2 := models.NewMember(uuid.New(), uuid.New(), uuid.New(), uuid.NewString(), true)
	member3 := models.NewMember(uuid.New(), uuid.New(), member1.GroupID(), uuid.NewString(), false)
	message1 := models.NewMessage(uuid.New(), member1.GroupID(), member1.ID(), uuid.NewString(), time.Now(), nil)
	message2 := models.NewMessage(uuid.New(), member2.GroupID(), member2.ID(), uuid.NewString(), time.Now(), nil)

	s.Require().NoError(s.repo.CreateMember(ctx, member1))
	s.Require().NoError(s.repo.CreateMember(ctx, member2))
	s.Require().NoError(s.repo.CreateMember(ctx, member3))
	s.Require().NoError(s.repo.CreateMessage(ctx, message1))
	s.Require().NoError(s.repo.CreateMessage(ctx, message2))

	testCases := []struct {
		desc           string
		userID         uuid.UUID
		messageID      uuid.UUID
		expectedResult error
	}{
		{
			desc:           "message_not_found",
			userID:         member1.UserID(),
			messageID:      member1.ID(),
			expectedResult: apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", member1.ID())),
		},
		{
			desc:           "user_not_in_group",
			userID:         member2.UserID(),
			messageID:      message1.ID(),
			expectedResult: apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", message1.ID())),
		},
		{
			desc:           "user_without_rights",
			userID:         member3.UserID(),
			messageID:      message1.ID(),
			expectedResult: apperrors.NewForbidden("user has no right to delete message"),
		},
		{
			desc:           "success",
			userID:         member2.UserID(),
			messageID:      message2.ID(),
			expectedResult: nil,
		},
	}

	for _, tC := range testCases {
		s.T().Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.expectedResult, s.repo.DeleteMessageForEveryone(ctx, tC.userID, tC.messageID))
		})
	}
}

func (s *MessagesSuite) TestDeleteGroup() {
	s.T().Parallel()
	ctx := context.Background()

	member1 := models.NewMember(uuid.New(), uuid.New(), uuid.New(), uuid.NewString(), true)
	member2 := models.NewMember(uuid.New(), uuid.New(), member1.GroupID(), uuid.NewString(), false)
	message1 := models.NewMessage(uuid.New(), member1.GroupID(), member1.ID(), uuid.NewString(), time.Now(), nil)
	file1 := models.NewMessageFile(message1.ID(), uuid.NewString(), "png")
	file2 := models.NewMessageFile(message1.ID(), uuid.NewString(), "png")
	message1.AddFiles(file1, file2)
	message2 := models.NewMessage(uuid.New(), member1.GroupID(), member1.ID(), uuid.NewString(), time.Now(), nil)

	s.Require().NoError(s.repo.CreateMember(ctx, member1))
	s.Require().NoError(s.repo.CreateMember(ctx, member2))
	s.Require().NoError(s.repo.CreateMessage(ctx, message1))
	s.Require().NoError(s.repo.CreateMessage(ctx, message2))

	msgs, err := s.repo.GetGroupMessages(ctx, member1.UserID(), member1.GroupID(), 0, 4)
	s.Require().NoError(err)
	s.Len(msgs, 2)

	s.NoError(s.repo.DeleteGroup(ctx, member1.GroupID()))
	msgs, err = s.repo.GetGroupMessages(ctx, member1.UserID(), member1.GroupID(), 0, 4)
	s.Require().Error(err)
	s.Len(msgs, 0)
}
