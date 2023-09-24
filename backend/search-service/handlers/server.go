package handlers

import (
	"crypto/rsa"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/database"
)

type Server struct {
	DB        database.DBLayer
	Listener  msgqueue.EventListener
	PublicKey *rsa.PublicKey
}
