package handlers

import (
	"crypto/rsa"
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/gin-gonic/gin"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/storage"
)

const MAX_BODY_BYTES = 4194304

type Server struct {
	logger       *log.Logger
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

func (srv *Server) WithLogger(logger *log.Logger) *Server {
	srv.logger = logger
	return srv
}

func (srv *Server) HandleError(c *gin.Context, err error) {
	if srv.logger != nil {
		log.Println(err)
	}

	c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
	return
}
