package handlers

import (
	"crypto/rsa"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/storage"
)

const MAX_BODY_BYTES = 4194304

type Server struct {
	DB           database.DBLayer
	Storage      storage.StorageLayer
	PublicKey    *rsa.PublicKey
	MaxBodyBytes int64
	Emitter      msgqueue.EventEmiter
}

func NewServer(db database.DBLayer, storage storage.StorageLayer, pubKey *rsa.PublicKey, emiter msgqueue.EventEmiter) *Server {
	return &Server{
		DB:           db,
		Storage:      storage,
		MaxBodyBytes: MAX_BODY_BYTES,
		PublicKey:    pubKey,
		Emitter:      emiter,
	}
}
