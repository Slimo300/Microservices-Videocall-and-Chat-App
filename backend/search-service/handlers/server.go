package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/database"
)

type Server struct {
	DB          database.DBLayer
	Listener    msgqueue.EventListener
	AuthService auth.TokenClient
}
