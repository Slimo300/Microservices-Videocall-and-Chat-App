package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
)

type Server struct {
	DB           database.DBLayer
	TokenService auth.TokenClient
	Emitter      msgqueue.EventEmitter
	Listener     msgqueue.EventListener
}
