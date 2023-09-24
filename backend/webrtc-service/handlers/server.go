package handlers

import (
	"crypto/rsa"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/webrtc"
)

type Server struct {
	PublicKey *rsa.PublicKey
	DB        database.DBLayer
	Relay     *webrtc.RoomsRelay
}
