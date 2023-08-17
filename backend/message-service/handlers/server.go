package handlers

import (
	"crypto/rsa"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/storage"
)

type Server struct {
	DB        database.DBLayer
	PublicKey *rsa.PublicKey
	Emitter   msgqueue.EventEmiter
	Storage   storage.StorageLayer
}

func NewServer(db database.DBLayer, pubKey *rsa.PublicKey, emitter msgqueue.EventEmiter, storage storage.StorageLayer) *Server {
	return &Server{
		DB:        db,
		PublicKey: pubKey,
		Emitter:   emitter,
		Storage:   storage,
	}
}
