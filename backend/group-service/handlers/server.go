package handlers

import (
	"crypto/rsa"
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/gin-gonic/gin"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/service"
)

const MAX_BODY_BYTES = 4194304

type Server struct {
	logger       *log.Logger
	PublicKey    *rsa.PublicKey
	MaxBodyBytes int64
	Service      service.ServiceLayer
}

func NewServer(service service.ServiceLayer, pubKey *rsa.PublicKey) *Server {
	return &Server{
		Service:      service,
		MaxBodyBytes: MAX_BODY_BYTES,
		PublicKey:    pubKey,
	}
}

func (srv *Server) WithLogger(logger *log.Logger) *Server {
	srv.logger = logger
	return srv
}

func (srv *Server) HandleError(c *gin.Context, err error) {
	if srv.logger != nil {
		srv.logger.Println(err)
	}

	c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
}
