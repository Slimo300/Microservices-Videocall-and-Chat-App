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
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"

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

	err = s.repo.UpdateGroup(ctx, group.ID(), func(g *models.Group) error {
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

func (s *GroupsSuite) TestRespondInvite_InviteDeclined() {

	s.T().Parallel()

	ctx := context.Background()
	var err error

	creator := models.NewUser(uuid.New(), uuid.NewString())
	target := models.NewUser(uuid.New(), uuid.NewString())
	group := models.CreateGroup(creator.ID(), uuid.NewString())
	invite := models.CreateInvite(creator.ID(), target.ID(), group.ID())

	s.Require().NoError(s.repo.CreateUser(ctx, creator))
	s.Require().NoError(s.repo.CreateUser(ctx, target))
	_, err = s.repo.CreateGroup(ctx, group)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, invite)
	s.Require().NoError(err)

	s.Require().NoError(s.repo.UpdateInvite(ctx, invite.ID(), func(i *models.Invite) (*models.Member, error) {
		if err := i.AnswerInvite(target.ID(), false); err != nil {
			return nil, err
		}
		return nil, nil
	}))

	invite, err = s.repo.GetInviteByID(ctx, target.ID(), invite.ID())
	s.Require().NoError(err)
	s.Equal(models.INVITE_DECLINED, invite.Status())

	group, err = s.repo.GetGroupByID(ctx, creator.ID(), invite.GroupID())
	s.Require().NoError(err)
	s.Len(group.Members(), 1)
}

func (s *GroupsSuite) TestRespondInvite_InviteAccepted() {
	s.T().Parallel()

	ctx := context.Background()
	var err error

	creator := models.NewUser(uuid.New(), uuid.NewString())
	target := models.NewUser(uuid.New(), uuid.NewString())
	group := models.CreateGroup(creator.ID(), uuid.NewString())
	invite := models.CreateInvite(creator.ID(), target.ID(), group.ID())

	s.Require().NoError(s.repo.CreateUser(ctx, creator))
	s.Require().NoError(s.repo.CreateUser(ctx, target))
	_, err = s.repo.CreateGroup(ctx, group)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, invite)
	s.Require().NoError(err)

	s.Require().NoError(s.repo.UpdateInvite(ctx, invite.ID(), func(i *models.Invite) (*models.Member, error) {
		if err := i.AnswerInvite(target.ID(), true); err != nil {
			return nil, err
		}
		member := models.NewMember(target.ID(), group.ID())
		return &member, nil
	}))

	invite, err = s.repo.GetInviteByID(ctx, target.ID(), invite.ID())
	s.Require().NoError(err)
	s.Equal(models.INVITE_ACCEPTED, invite.Status())
	group, err = s.repo.GetGroupByID(ctx, creator.ID(), group.ID())
	s.Require().NoError(err)
	s.Len(group.Members(), 2)
	member, _ := group.GetMemberByUserID(invite.TargetID())
	s.Equal(group.ID(), member.GroupID())
	s.Equal(target.ID(), member.UserID())
}

func (s *GroupsSuite) TestCreateInvite() {
	s.T().Parallel()

	ctx := context.Background()
	var err error

	creator := models.NewUser(uuid.New(), uuid.NewString())
	target := models.NewUser(uuid.New(), uuid.NewString())
	userNotInGroup := models.NewUser(uuid.New(), uuid.NewString())
	userInGroup := models.NewUser(uuid.New(), uuid.NewString())
	userInvited := models.NewUser(uuid.New(), uuid.NewString())
	group := models.CreateGroup(creator.ID(), uuid.NewString())
	inviteAnswered := models.CreateInvite(creator.ID(), userInGroup.ID(), group.ID())
	inviteAwaiting := models.CreateInvite(creator.ID(), userInvited.ID(), group.ID())

	s.Require().NoError(s.repo.CreateUser(ctx, creator))
	s.Require().NoError(s.repo.CreateUser(ctx, target))
	s.Require().NoError(s.repo.CreateUser(ctx, userInGroup))
	s.Require().NoError(s.repo.CreateUser(ctx, userInvited))
	_, err = s.repo.CreateGroup(ctx, group)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, inviteAnswered)
	s.Require().NoError(err)
	_, err = s.repo.CreateInvite(ctx, inviteAwaiting)
	s.Require().NoError(err)
	s.Require().NoError(s.repo.UpdateInvite(ctx, inviteAnswered.ID(), func(i *models.Invite) (*models.Member, error) {
		if err := i.AnswerInvite(userInGroup.ID(), true); err != nil {
			return nil, err
		}
		member := models.NewMember(userInGroup.ID(), group.ID())
		return &member, nil
	}))

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
