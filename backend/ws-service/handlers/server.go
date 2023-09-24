package handlers

import (
	"crypto/rsa"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/ws"
)

type Server struct {
	DB        database.DBLayer
	Hub       *ws.WSHub
	PublicKey *rsa.PublicKey
}
