package handlers

import (
	"crypto/rsa"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/storage"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/email"
)

type Server struct {
	DB           database.DBLayer
	Emitter      msgqueue.EventEmiter
	ImageStorage storage.StorageLayer
	TokenKey     *rsa.PublicKey
	TokenClient  auth.TokenServiceClient
	EmailClient  email.EmailServiceClient
	MaxBodyBytes int64
	Domain       string
}
