package database_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database/orm"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	merrors "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models/errors"
)

var pool *dockerTestPool

type dockerTestPool struct {
	*sync.Mutex
	pool      *dockertest.Pool
	resources []*dockertest.Resource
}

func (p *dockerTestPool) createRepository(image, tag string, env []string, connect func(*dockertest.Resource) (database.GroupsRepository, error)) database.GroupsRepository {
	resource, err := p.pool.Run(image, tag, env)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	var repo database.GroupsRepository
	if err := p.pool.Retry(func() error {
		repo, err = connect(resource)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	p.addResource(resource)
	return repo
}

func (p *dockerTestPool) addResource(resource *dockertest.Resource) {
	p.Lock()
	defer p.Unlock()

	p.resources = append(p.resources, resource)
}

func (p *dockerTestPool) clear() {
	p.Lock()
	defer p.Unlock()

	for _, resource := range p.resources {
		if err := p.pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}

func newPool() (*dockerTestPool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return &dockerTestPool{}, fmt.Errorf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		return &dockerTestPool{}, fmt.Errorf("Could not connect to Docker: %s", err)
	}
	return &dockerTestPool{
		Mutex: &sync.Mutex{},
		pool:  pool,
	}, nil
}

func TestMain(m *testing.M) {
	var err error
	pool, err = newPool()
	if err != nil {
		log.Fatalf("Could not construct pool: %v", err)
	}

	code := m.Run()

	pool.clear()
	os.Exit(code)
}

type GroupsSuite struct {
	suite.Suite
	repo database.GroupsRepository
}

func TestSuites(t *testing.T) {
	gormRepository := pool.createRepository("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"}, func(r *dockertest.Resource) (database.GroupsRepository, error) {
		return orm.NewGroupsGormRepository(fmt.Sprintf("root:secret@(localhost:%s)/mysql", r.GetPort("3306/tcp")))
	})
	for _, s := range []*GroupsSuite{{repo: gormRepository}} {
		s := s
		suite.Run(t, s)
	}
}

func (s *GroupsSuite) TestUsers() {
	s.T().Parallel()

	ctx := context.Background()

	user := models.NewUser(uuid.New(), uuid.NewString())
	err := s.repo.CreateUser(ctx, user)
	s.Require().NoError(err)

	returnedUser, err := s.repo.GetUserByID(ctx, user.ID())
	s.Require().NoError(err)
	s.Equal(user.ID(), returnedUser.ID())
	s.Equal(user.Username(), returnedUser.Username())
	s.Equal(false, returnedUser.HasPicture())

	err = s.repo.UpdateUser(ctx, user.ID(), func(u *models.User) error {
		u.UpdatePictureState(true)
		return nil
	})
	s.Require().NoError(err)

	returnedUser, err = s.repo.GetUserByID(ctx, user.ID())
	s.Require().NoError(err)
	s.Equal(user.Username(), returnedUser.Username())
	s.Equal(true, returnedUser.HasPicture())
}

func (s *GroupsSuite) TestGroups() {
	s.T().Parallel()

	ctx := context.Background()

	user := models.NewUser(uuid.New(), uuid.NewString())
	err := s.repo.CreateUser(ctx, user)
	s.Require().NoError(err)

	group := models.CreateGroup(user.ID(), uuid.NewString())
	group, err = s.repo.CreateGroup(ctx, group)
	s.Require().NoError(err)
	s.Equal(user.Username(), group.Members()[0].User().Username())

	returnedGroup, err := s.repo.GetGroupByID(ctx, user.ID(), group.ID())
	s.Require().NoError(err)
	s.Equal(group.ID(), returnedGroup.ID())
	s.Equal(group.Name(), returnedGroup.Name())
	s.Len(group.Members(), 1)

	member := returnedGroup.Members()[0]
	s.Equal(user.ID(), member.UserID())
	s.Equal(user.Username(), member.User().Username())

	err = s.repo.UpdateGroup(ctx, user.ID(), group.ID(), func(g *models.Group) error {
		_ = g.ChangePictureStateIfIncorrect(true)
		return nil
	})
	s.Require().NoError(err)

	returnedGroup, err = s.repo.GetGroupByID(ctx, user.ID(), group.ID())
	s.Require().NoError(err)
	s.True(returnedGroup.HasPicture())

	err = s.repo.DeleteGroup(ctx, user.ID(), group.ID())
	s.Require().NoError(err)

	_, err = s.repo.GetGroupByID(ctx, user.ID(), group.ID())
	s.Error(err)
}

func (s *GroupsSuite) TestRespondInvite() {
	// s.T().Parallel()

	creator := models.NewUser(uuid.New(), uuid.NewString())
	target1 := models.NewUser(uuid.New(), uuid.NewString())
	target2 := models.NewUser(uuid.New(), uuid.NewString())
	group := models.CreateGroup(creator.ID(), uuid.NewString())
	invite1 := models.CreateInvite(creator.ID(), target1.ID(), group.ID())
	invite2 := models.CreateInvite(creator.ID(), target2.ID(), group.ID())

	ctx := context.Background()
	var err error
	s.Require().NoError(s.repo.CreateUser(ctx, creator))
	s.Require().NoError(s.repo.CreateUser(ctx, target1))
	s.Require().NoError(s.repo.CreateUser(ctx, target2))
	_, err = s.repo.CreateGroup(ctx, group)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, invite1)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, invite2)
	s.Require().NoError(err)

	testCases := []struct {
		desc          string
		inviteID      uuid.UUID
		targetID      uuid.UUID
		updateFn      func(*models.Invite) (*models.Member, error)
		expectedError error
		memberCreated bool
	}{
		{
			desc:          "invite_not_found",
			inviteID:      group.ID(),
			expectedError: merrors.NewInviteNotFoundError(group.ID().String()),
		},
		{
			desc:          "updateFn_error",
			inviteID:      invite1.ID(),
			expectedError: merrors.NewInviteAlreadyAnsweredError(invite1.ID().String()),
			updateFn: func(i *models.Invite) (*models.Member, error) {
				return nil, merrors.NewInviteAlreadyAnsweredError(invite1.ID().String())
			},
		},
		{
			desc:     "invite_declined",
			inviteID: invite1.ID(),
			targetID: target1.ID(),
			updateFn: func(i *models.Invite) (*models.Member, error) {
				return nil, i.AnswerInvite(target1.ID(), false)
			},
		},
		{
			desc:     "invite_accepted",
			inviteID: invite2.ID(),
			targetID: target2.ID(),
			updateFn: func(i *models.Invite) (*models.Member, error) {
				m := models.NewMember(target2.ID(), group.ID())
				return &m, i.AnswerInvite(target2.ID(), false)
			},
			memberCreated: true,
		},
	}
	for _, tC := range testCases {
		tC := tC
		s.T().Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			err := s.repo.UpdateInvite(ctx, tC.inviteID, tC.updateFn)
			require.Equal(t, tC.expectedError, err)

			if tC.expectedError != nil {
				return
			}
			invite, err := s.repo.GetInviteByID(ctx, tC.targetID, tC.inviteID)
			require.NoError(t, err)
			assert.NotEqual(t, models.INVITE_AWAITING, invite.Status())

			if !tC.memberCreated {
				return
			}

			group, err := s.repo.GetGroupByID(ctx, tC.targetID, invite.GroupID())
			require.NoError(t, err)
			member, _ := group.GetMemberByUserID(tC.targetID)
			assert.Equal(t, group.ID(), member.GroupID())
			assert.Equal(t, target2.ID(), member.UserID())
		})
	}
}

func (s *GroupsSuite) TestCreateInvite() {
	// s.T().Parallel()

	creator := models.NewUser(uuid.New(), uuid.NewString())
	target := models.NewUser(uuid.New(), uuid.NewString())
	userNotInGroup := models.NewUser(uuid.New(), uuid.NewString())
	userInGroup := models.NewUser(uuid.New(), uuid.NewString())
	userInvited := models.NewUser(uuid.New(), uuid.NewString())
	group := models.CreateGroup(creator.ID(), uuid.NewString())
	_ = group.AddMember(userInGroup.ID())
	inviteAwaiting := models.CreateInvite(creator.ID(), userInvited.ID(), group.ID())

	ctx := context.Background()
	var err error
	s.Require().NoError(s.repo.CreateUser(ctx, creator))
	s.Require().NoError(s.repo.CreateUser(ctx, target))
	s.Require().NoError(s.repo.CreateUser(ctx, userInGroup))
	s.Require().NoError(s.repo.CreateUser(ctx, userInvited))
	_, err = s.repo.CreateGroup(ctx, group)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, inviteAwaiting)
	s.Require().NoError(err)

	testCases := []struct {
		name           string
		issuerID       uuid.UUID
		targetID       uuid.UUID
		groupID        uuid.UUID
		expectedError  error
		expectedResult models.Invite
	}{
		{
			name:          "user_not_in_group",
			issuerID:      userNotInGroup.ID(),
			targetID:      target.ID(),
			groupID:       group.ID(),
			expectedError: merrors.NewUserNotInGroupError(userNotInGroup.ID().String(), group.ID().String()),
		},
		{
			name:          "user_without_rights",
			issuerID:      userInGroup.ID(),
			targetID:      target.ID(),
			groupID:       group.ID(),
			expectedError: merrors.NewMemberUnauthorizedError(group.ID().String(), merrors.AddMemberAction()),
		},
		{
			name:          "user_not_found",
			issuerID:      creator.ID(),
			targetID:      group.ID(),
			groupID:       group.ID(),
			expectedError: merrors.NewUserNotFoundError(group.ID().String()),
		},
		{
			name:          "user_already_in_group",
			issuerID:      creator.ID(),
			targetID:      userInGroup.ID(),
			groupID:       group.ID(),
			expectedError: merrors.NewUserAlreadyInGroupError(userInGroup.ID().String(), group.ID().String()),
		},
		{
			name:          "user_already_invited",
			issuerID:      creator.ID(),
			targetID:      userInvited.ID(),
			groupID:       group.ID(),
			expectedError: merrors.NewUserAlreadyInvitedError(userInvited.ID().String(), group.ID().String()),
		},
		{
			name:          "invite_created",
			issuerID:      creator.ID(),
			targetID:      target.ID(),
			groupID:       group.ID(),
			expectedError: nil,
		},
	}

	for i := range testCases {
		tC := testCases[i]
		s.T().Run(tC.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			invite := models.CreateInvite(tC.issuerID, tC.targetID, tC.groupID)

			_, err := s.repo.CreateInvite(ctx, invite)
			assert.Equal(t, tC.expectedError, err)
		})
	}
}

func (s *GroupsSuite) TestDeleteMember() {
	// s.T().Parallel()

	creator := models.NewUser(uuid.New(), uuid.NewString())
	userNoRights := models.NewUser(uuid.New(), uuid.NewString())
	userToDelete := models.NewUser(uuid.New(), uuid.NewString())
	userToFailDelete := models.NewUser(uuid.New(), uuid.NewString())
	group := models.CreateGroup(creator.ID(), uuid.NewString())
	_ = group.AddMember(userNoRights.ID())
	memberToFailDelete := group.AddMember(userToFailDelete.ID())
	memberToDelete := group.AddMember(userToDelete.ID())

	ctx := context.Background()

	s.Require().NoError(s.repo.CreateUser(ctx, creator))
	s.Require().NoError(s.repo.CreateUser(ctx, userNoRights))
	s.Require().NoError(s.repo.CreateUser(ctx, userToDelete))
	s.Require().NoError(s.repo.CreateUser(ctx, userToFailDelete))
	group, err := s.repo.CreateGroup(ctx, group)
	s.Require().NoError(err)

	testCases := []struct {
		desc           string
		issuerID       uuid.UUID
		memberID       uuid.UUID
		expectedResult error
	}{
		{
			desc:           "member_not_found",
			memberID:       group.ID(),
			issuerID:       creator.ID(),
			expectedResult: merrors.NewMemberNotFoundError(group.ID().String()),
		},
		{
			desc:           "issuer_not_in_group",
			memberID:       memberToFailDelete.ID(),
			issuerID:       group.ID(),
			expectedResult: merrors.NewUserNotInGroupError(group.ID().String(), group.ID().String()),
		},
		{
			desc:           "member_cant_delete",
			memberID:       memberToFailDelete.ID(),
			issuerID:       userNoRights.ID(),
			expectedResult: merrors.NewMemberUnauthorizedError(group.ID().String(), merrors.DeleteMemberAction()),
		},
		{
			desc:           "success",
			memberID:       memberToDelete.ID(),
			issuerID:       creator.ID(),
			expectedResult: nil,
		},
	}

	for _, tC := range testCases {
		tC := tC
		s.T().Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			err := s.repo.DeleteMember(ctx, tC.issuerID, tC.memberID)
			assert.Equal(t, tC.expectedResult, err)
		})
	}
}

func (s *GroupsSuite) TestUpdateMember() {
	// s.T().Parallel()

	creator := models.NewUser(uuid.New(), uuid.NewString())
	userToUpdate := models.NewUser(uuid.New(), uuid.NewString())
	userNoRights := models.NewUser(uuid.New(), uuid.NewString())

	group := models.CreateGroup(creator.ID(), uuid.NewString())
	memberToUpdate := group.AddMember(userToUpdate.ID())
	_ = group.AddMember(userNoRights.ID())

	ctx := context.Background()
	s.Require().NoError(s.repo.CreateUser(ctx, creator))
	s.Require().NoError(s.repo.CreateUser(ctx, userToUpdate))
	s.Require().NoError(s.repo.CreateUser(ctx, userNoRights))

	_, err := s.repo.CreateGroup(ctx, group)
	s.Require().NoError(err)

	testCases := []struct {
		desc           string
		memberID       uuid.UUID
		issuerID       uuid.UUID
		updateFn       func(m *models.Member) error
		expectedResult error
	}{
		{
			desc:           "member_not_found",
			issuerID:       creator.ID(),
			memberID:       group.ID(),
			expectedResult: merrors.NewMemberNotFoundError(group.ID().String()),
		},
		{
			desc:           "issuer_not_in_group",
			issuerID:       group.ID(),
			memberID:       memberToUpdate.ID(),
			expectedResult: merrors.NewUserNotInGroupError(group.ID().String(), group.ID().String()),
		},
		{
			desc:           "issuer_without_rights",
			issuerID:       userNoRights.ID(),
			memberID:       memberToUpdate.ID(),
			expectedResult: merrors.NewMemberUnauthorizedError(group.ID().String(), merrors.UpdateMemberAction()),
		},
		{
			desc:     "updateFn_error",
			issuerID: creator.ID(),
			memberID: memberToUpdate.ID(),
			updateFn: func(m *models.Member) error {
				return merrors.NewMemberUnauthorizedError(group.ID().String(), merrors.UpdateMemberAction())
			},
			expectedResult: merrors.NewMemberUnauthorizedError(group.ID().String(), merrors.UpdateMemberAction()),
		},
		{
			desc:           "success",
			issuerID:       creator.ID(),
			memberID:       memberToUpdate.ID(),
			updateFn:       func(m *models.Member) error { m.ApplyRights(models.MemberRights{Adding: true}); return nil },
			expectedResult: nil,
		},
	}

	for _, tC := range testCases {
		tC := tC
		s.T().Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			err := s.repo.UpdateMember(ctx, tC.issuerID, tC.memberID, tC.updateFn)
			assert.Equal(t, tC.expectedResult, err)
		})
	}
}

func (s *GroupsSuite) TestGetUserInvites() {
	s.T().Parallel()

	user1 := models.NewUser(uuid.New(), uuid.NewString())
	user2 := models.NewUser(uuid.New(), uuid.NewString())
	user3 := models.NewUser(uuid.New(), uuid.NewString())

	group1 := models.CreateGroup(user1.ID(), uuid.NewString())
	group2 := models.CreateGroup(user2.ID(), uuid.NewString())

	invite1 := models.CreateInvite(user1.ID(), user2.ID(), group1.ID())
	invite2 := models.CreateInvite(user2.ID(), user1.ID(), group2.ID())
	invite3 := models.CreateInvite(user2.ID(), user3.ID(), group2.ID())

	ctx := context.Background()
	s.Require().NoError(s.repo.CreateUser(ctx, user1))
	s.Require().NoError(s.repo.CreateUser(ctx, user2))
	s.Require().NoError(s.repo.CreateUser(ctx, user3))
	_, err := s.repo.CreateGroup(ctx, group1)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, invite1)
	s.Require().NoError(err)
	_, err = s.repo.CreateGroup(ctx, group2)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, invite2)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, invite3)
	s.Require().NoError(err)

	invites, err := s.repo.GetUserInvites(ctx, user1.ID(), 4, 0)
	s.Require().NoError(err)
	s.True(findInvites(invites, invite1.ID()), invite2.ID())
	s.False(findInvites(invites, invite3.ID()))
}

func findInvites(invites []models.Invite, inviteIDs ...uuid.UUID) bool {
	if invites == nil {
		return false
	}
	searchedInvites := make(map[uuid.UUID]bool)
	for i := range inviteIDs {
		searchedInvites[inviteIDs[i]] = true
	}
	for i := range invites {
		delete(searchedInvites, invites[i].ID())
	}
	return len(searchedInvites) == 0
}

func (s *GroupsSuite) TestGetUserGroups() {
	s.T().Parallel()

	user1 := models.NewUser(uuid.New(), uuid.NewString())
	user2 := models.NewUser(uuid.New(), uuid.NewString())

	group1 := models.CreateGroup(user1.ID(), uuid.NewString())
	group2 := models.CreateGroup(user2.ID(), uuid.NewString())
	group2.AddMember(user1.ID())
	group3 := models.CreateGroup(user2.ID(), uuid.NewString())

	ctx := context.Background()
	s.Require().NoError(s.repo.CreateUser(ctx, user1))
	s.Require().NoError(s.repo.CreateUser(ctx, user2))
	_, err := s.repo.CreateGroup(ctx, group1)
	s.Require().NoError(err)
	_, err = s.repo.CreateGroup(ctx, group2)
	s.Require().NoError(err)
	_, err = s.repo.CreateGroup(ctx, group3)
	s.Require().NoError(err)

	groups, err := s.repo.GetUserGroups(ctx, user1.ID())
	s.Require().NoError(err)
	s.True(findGroups(groups, group1.ID(), group2.ID()))
	s.False(findGroups(groups, group3.ID()))
}

func findGroups(groups []models.Group, groupIDs ...uuid.UUID) bool {
	if groups == nil {
		return false
	}
	searchedGroups := make(map[uuid.UUID]bool)
	for i := range groupIDs {
		searchedGroups[groupIDs[i]] = true
	}
	for i := range groups {
		delete(searchedGroups, groups[i].ID())
	}
	return len(searchedGroups) == 0
}
