package handlers

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/webrtc"
)

type Server struct {
	TokenClient auth.TokenClient
	DB          database.DBLayer
	Rooms       map[string]*webrtc.Room
}
