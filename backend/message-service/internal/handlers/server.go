package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/chat-messageservice/internal/database"
	"github.com/Slimo300/chat-messageservice/internal/storage"
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
)

type Server struct {
	DB          database.DBLayer
	TokenClient tokens.TokenClient
	Emitter     msgqueue.EventEmiter
	Storage     storage.StorageLayer
}

func NewServer(db database.DBLayer, tokenClient tokens.TokenClient, emitter msgqueue.EventEmiter, storage storage.StorageLayer) *Server {
	return &Server{
		DB:          db,
		TokenClient: tokenClient,
		Emitter:     emitter,
		Storage:     storage,
	}
}
