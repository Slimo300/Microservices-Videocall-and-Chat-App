package database_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database/orm"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
)

var pool *dockerTestPool

type dockerTestPool struct {
	*sync.Mutex
	pool      *dockertest.Pool
	resources []*dockertest.Resource
}

func (p *dockerTestPool) createRepository(image, tag string, env []string, connect func(*dockertest.Resource) (database.UsersRepository, error)) database.UsersRepository {
	resource, err := p.pool.Run(image, tag, env)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	var repo database.UsersRepository
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

func TestSuites(t *testing.T) {
	gormRepository := pool.createRepository("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"}, func(r *dockertest.Resource) (database.UsersRepository, error) {
		return orm.NewUsersGormRepository(fmt.Sprintf("root:secret@(localhost:%s)/mysql", r.GetPort("3306/tcp")))
	})
	for _, s := range []*UserSuite{{repo: gormRepository}} {
		s := s
		suite.Run(t, s)
	}
}

type UserSuite struct {
	suite.Suite
	repo database.UsersRepository
}

func (s *UserSuite) TestRegisterAndDelete() {
	s.T().Parallel()
	if s.repo == nil {
		panic("repo is nil")
	}
	ctx := context.Background()

	returnedUser, err := s.repo.GetUserByEmail(ctx, "host@net.pl")
	s.Require().Error(err)
	s.Require().Nil(returnedUser)

	user := models.MustNewUser("host@net.pl", "username", "password")
	code := models.NewAuthorizationCode(user.ID(), models.EmailVerificationCode)

	err = s.repo.RegisterUser(ctx, user, code)
	s.Require().Nil(err)

	returnedUser, err = s.repo.GetUserByEmail(ctx, "host@net.pl")
	s.Require().NoError(err)
	s.Equal(user.Email(), returnedUser.Email())
	s.Equal(user.PasswordHash(), returnedUser.PasswordHash())
	s.Equal(user.Username(), returnedUser.Username())

	err = s.repo.DeleteUser(ctx, user.ID())
	s.Require().NoError(err)
	_, err = s.repo.GetUserByEmail(ctx, "host@net.pl")
	s.Error(err)
}

func (s *UserSuite) TestUpdateUSerByID() {
	s.T().Parallel()
	if s.repo == nil {
		panic("repo is nil")
	}
	ctx := context.Background()
	user := models.MustNewUser("host1@net.pl", "username1", "password")
	code := models.NewAuthorizationCode(user.ID(), models.EmailVerificationCode)

	_, err := s.repo.GetUserByID(ctx, user.ID())
	s.Require().Error(err)

	err = s.repo.RegisterUser(ctx, user, code)
	s.Require().NoError(err)
	returnedUser, err := s.repo.GetUserByID(ctx, user.ID())
	s.Require().NoError(err)
	s.Require().Equal(returnedUser.PasswordHash(), user.PasswordHash())

	err = s.repo.UpdateUserByID(ctx, user.ID(), func(u *models.User) (*models.User, error) {
		_ = u.SetPassword("newPassword")
		return u, errors.New("error")
	})
	s.Error(err)

	err = s.repo.UpdateUserByID(ctx, user.ID(), func(u *models.User) (*models.User, error) {
		_ = u.SetPassword("newPassword")
		return u, nil
	})
	s.NoError(err)
	returnedUser, err = s.repo.GetUserByID(ctx, user.ID())
	s.NoError(err)
	s.Equal(user.Email(), returnedUser.Email())
}

func (s *UserSuite) TestUpdateUserByCode() {
	s.T().Parallel()
	if s.repo == nil {
		panic("repo is nil")
	}
	ctx := context.Background()

	user := models.MustNewUser("host2@net.pl", "username2", "password")
	code := models.NewAuthorizationCode(user.ID(), models.EmailVerificationCode)

	err := s.repo.RegisterUser(ctx, user, code)
	s.Require().NoError(err)
	_, err = s.repo.GetUserByID(ctx, user.ID())
	s.Require().NoError(err)

	err = s.repo.UpdateUserByCode(ctx, code.Code(), models.EmailVerificationCode, func(u *models.User) (*models.User, error) {
		u.Verify()
		return nil, errors.New("error")
	})
	s.Require().Error(err)

	err = s.repo.UpdateUserByCode(ctx, code.Code(), models.ResetPasswordCode, func(u *models.User) (*models.User, error) {
		u.Verify()
		return u, nil
	})
	s.Require().Error(err)

	err = s.repo.UpdateUserByCode(ctx, code.Code(), models.EmailVerificationCode, func(u *models.User) (*models.User, error) {
		u.Verify()
		return u, nil
	})
	s.Require().NoError(err)
	user, err = s.repo.GetUserByID(ctx, user.ID())
	s.Require().NoError(err)
	s.True(user.Verified())

	code = models.NewAuthorizationCode(user.ID(), models.ResetPasswordCode)
	err = s.repo.CreateAuthorizationCode(ctx, code)
	s.Require().NoError(err)

	err = s.repo.UpdateUserByCode(ctx, code.Code(), models.ResetPasswordCode, func(u *models.User) (*models.User, error) {
		_ = u.ChangePictureStateIfIncorrect(true)
		return u, nil
	})
	s.Require().NoError(err)
	user, err = s.repo.GetUserByID(ctx, user.ID())
	s.NoError(err)
	s.True(user.HasPicture())
}
