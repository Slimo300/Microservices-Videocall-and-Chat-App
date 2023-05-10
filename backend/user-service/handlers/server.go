package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/storage"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/email"
)

type Server struct {
	DB           database.DBLayer
	Emitter      msgqueue.EventEmiter
	ImageStorage storage.StorageLayer
	TokenClient  auth.TokenClient
	EmailClient  email.EmailClient
	MaxBodyBytes int64
	Domain       string
}
