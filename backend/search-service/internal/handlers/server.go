package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/chat-searchservice/internal/database"

	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
)

type Server struct {
	DB          database.DBLayer
	Listener    msgqueue.EventListener
	TokenClient tokens.TokenClient
}
