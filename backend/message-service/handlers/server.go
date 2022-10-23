package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
)

type Server struct {
	DB           database.DBLayer
	TokenService auth.TokenClient
	domain       string
}
