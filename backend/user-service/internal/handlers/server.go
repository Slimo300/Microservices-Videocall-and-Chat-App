package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"

	emails "github.com/Slimo300/chat-emailservice/pkg/client"
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"

	"github.com/Slimo300/chat-userservice/internal/database"
	"github.com/Slimo300/chat-userservice/internal/storage"
)

type Server struct {
	DB           database.DBLayer
	Emitter      msgqueue.EventEmiter
	ImageStorage storage.StorageLayer
	TokenClient  tokens.TokenClient
	EmailClient  emails.EmailClient
	MaxBodyBytes int64
	Domain       string
}
