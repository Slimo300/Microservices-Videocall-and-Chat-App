package handlers

import (
	"crypto/rsa"
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app"
)

type Server struct {
	app          app.App
	publicKey    *rsa.PublicKey
	maxBodyBytes int64
	domain       string
}

func NewServer(app app.App, publicKey *rsa.PublicKey, domain string, maxBodyBytes int64) http.Handler {
	server := Server{
		app:          app,
		publicKey:    publicKey,
		domain:       domain,
		maxBodyBytes: maxBodyBytes,
	}
	return server.setup(domain)
}

func NewTestServer(app app.App) Server {
	return Server{
		app: app,
	}
}
