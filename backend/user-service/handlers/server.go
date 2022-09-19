package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/storage"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/email"
)

type Server struct {
	DB           database.DBlayer
	Storage      storage.StorageLayer
	TokenService auth.TokenClient
	EmailService email.EmailService
	domain       string
	MaxBodyBytes int64
}
