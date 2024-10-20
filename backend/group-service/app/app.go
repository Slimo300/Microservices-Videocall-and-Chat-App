package app

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app/command"
	query "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app/query"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"

	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database/mock"
	mockqueue "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/mock"
	mockstorage "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage/mock"
)

type Commands struct {
	CreateGroup        command.CreateGroupHandler
	DeleteGroup        command.DeleteGroupHandler
	SetGroupPicture    command.SetGroupPictureHandler
	DeleteGroupPicture command.DeleteGroupPictureHandler

	GrantRights  command.GrantRightsHandler
	DeleteMember command.DeleteMemberHandler

	SendInvite    command.SendInviteHandler
	RespondInvite command.RespondInviteHandler

	CreateUser command.CreateUserHandler
	UpdateUser command.UpdateUserHandler
}

type Queries struct {
	GetUserInvites query.GetUserInvitesHandler
	GetUserGroups  query.GetUserGroupsHandler
}

type App struct {
	Commands Commands
	Queries  Queries
}

func NewTestApplication() App {
	repo := new(mockdb.GroupsMockRepository)
	emitter := new(mockqueue.MockEmitter)
	storage := new(mockstorage.MockStorage)

	return App{
		Queries: Queries{
			GetUserInvites: query.NewGetUserInvitesHandler(repo),
			GetUserGroups:  query.NewGetUserGroupsHandler(repo),
		},
		Commands: Commands{
			CreateGroup:        command.NewCreateGroupHandler(repo, emitter),
			DeleteGroup:        command.NewDeleteGroupHandler(repo, emitter),
			SetGroupPicture:    command.NewSetGroupPictureHandler(repo, emitter, storage),
			DeleteGroupPicture: command.NewDeleteGroupPictureHandler(repo, emitter, storage),
			DeleteMember:       command.NewDeleteMemberHandler(repo, emitter),
			GrantRights:        command.NewGrantRightsHandler(repo, emitter),
			SendInvite:         command.NewSendInviteHandler(repo, emitter),
			RespondInvite:      command.NewRespondInviteHandler(repo, emitter),
		},
	}
}

func NewApplication(repo database.GroupsRepository, storage storage.Storage, emitter msgqueue.EventEmiter) App {
	return App{
		Queries: Queries{
			GetUserInvites: query.NewGetUserInvitesHandler(repo),
			GetUserGroups:  query.NewGetUserGroupsHandler(repo),
		},
		Commands: Commands{
			CreateGroup:        command.NewCreateGroupHandler(repo, emitter),
			DeleteGroup:        command.NewDeleteGroupHandler(repo, emitter),
			SetGroupPicture:    command.NewSetGroupPictureHandler(repo, emitter, storage),
			DeleteGroupPicture: command.NewDeleteGroupPictureHandler(repo, emitter, storage),
			DeleteMember:       command.NewDeleteMemberHandler(repo, emitter),
			GrantRights:        command.NewGrantRightsHandler(repo, emitter),
			SendInvite:         command.NewSendInviteHandler(repo, emitter),
			RespondInvite:      command.NewRespondInviteHandler(repo, emitter),
			CreateUser:         command.NewCreateUserHandler(repo),
			UpdateUser:         command.NewUpdateUserHandler(repo),
		},
	}
}
