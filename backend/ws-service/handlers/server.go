package handlers

import (
	"crypto/rsa"

	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/ws"
)

type Server struct {
	DB        database.DBLayer
	Hub       *ws.WSHub
	PublicKey *rsa.PublicKey
}
