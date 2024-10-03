package handlers

import (
	"crypto/rsa"
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/gin-gonic/gin"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/database"
)

type Server struct {
	DB        database.DBLayer
	Listener  msgqueue.EventListener
	PublicKey *rsa.PublicKey
}

func NewServer(db database.DBLayer, pubKey *rsa.PublicKey, origin string) http.Handler {
	server := Server{
		DB:        db,
		PublicKey: pubKey,
	}

	return server.setup(origin)
}

func (s *Server) setup(origin string) http.Handler {
	engine := gin.Default()

	engine.Use(corsMiddleware(origin))
	engine.Use(auth.MustAuthWithKey(s.PublicKey))

	engine.GET("/search/:name", s.SearchUsers)

	return engine
}

func corsMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
