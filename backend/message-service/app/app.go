package app

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/app/command"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/app/query"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
)

type Commands struct {
	CreateMember             command.CreateMemberHandler
	UpdateMember             command.UpdateMemberHandler
	DeleteMember             command.DeleteMemberHandler
	DeleteGroup              command.DeleteGroupHandler
	CreateMessage            command.CreateMessageHandler
	DeleteMessageForYourself command.DeleteMessageForYourselfHandler
	DeleteMessageForEveryone command.DeleteMessageForEveryoneHandler
}

type Queries struct {
	GetGroupMessages               query.GetGroupMessagesHandler
	GetPresignedGetRequestsHandler query.GetPresignedGetRequestsHandler
	GetPresignedPutRequestsHandler query.GetPresignedPutRequestsHandler
}

type App struct {
	Commands Commands
	Queries  Queries
}

func NewApplication(repo database.MessagesRepository, storage storage.Storage, emitter msgqueue.EventEmiter) App {
	return App{
		Queries: Queries{
			GetGroupMessages:               query.NewGetGroupMessagesHandler(repo),
			GetPresignedGetRequestsHandler: query.NewGetPresignedGetRequestsHandler(repo, storage),
			GetPresignedPutRequestsHandler: query.NewGetPresignedPutRequestsHandler(repo, storage),
		},
		Commands: Commands{
			CreateMember:             command.NewCreateMemberHandler(repo),
			UpdateMember:             command.NewUpdateMemberHandler(repo),
			DeleteMember:             command.NewDeleteMemberHandler(repo),
			DeleteGroup:              command.NewDeleteGroupHandler(repo, storage),
			CreateMessage:            command.NewCreateMessageHandler(repo),
			DeleteMessageForYourself: command.NewDeleteMessageForYourselfHandler(repo),
			DeleteMessageForEveryone: command.NewDeleteMessageForEveryoneHandler(repo, storage, emitter),
		},
	}
}
