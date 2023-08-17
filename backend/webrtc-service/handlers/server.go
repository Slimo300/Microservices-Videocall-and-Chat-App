package handlers

import (
	"crypto/rsa"

	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/webrtc"
)

type Server struct {
	PublicKey *rsa.PublicKey
	DB        database.DBLayer
	Relay     *webrtc.RoomsRelay
}
