package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
)

type Server struct {
	DB           database.DBLayer
	TokenService auth.TokenClient
	Emitter      msgqueue.EventEmiter
	Listener     msgqueue.EventListener
	Storage      storage.FileStorage
}

func NewServer(db database.DBLayer, auth auth.TokenClient, emitter msgqueue.EventEmiter, listener msgqueue.EventListener, storage storage.FileStorage) *Server {
	return &Server{
		DB:           db,
		TokenService: auth,
		Emitter:      emitter,
		Listener:     listener,
		Storage:      storage,
	}
}
