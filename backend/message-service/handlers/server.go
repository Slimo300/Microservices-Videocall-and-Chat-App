package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/storage"
)

type Server struct {
	DB          database.DBLayer
	TokenClient auth.TokenClient
	Emitter     msgqueue.EventEmiter
	Storage     storage.StorageLayer
}

func NewServer(db database.DBLayer, tokenClient auth.TokenClient, emitter msgqueue.EventEmiter, storage storage.StorageLayer) *Server {
	return &Server{
		DB:          db,
		TokenClient: tokenClient,
		Emitter:     emitter,
		Storage:     storage,
	}
}
