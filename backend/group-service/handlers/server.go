package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/storage"
	tokens "github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/gin-gonic/gin"
)

const MAX_BODY_BYTES = 4194304

type Server struct {
	DB           database.DBLayer
	Storage      storage.StorageLayer
	TokenClient  tokens.TokenClient
	MaxBodyBytes int64
	Emitter      msgqueue.EventEmiter
}

func NewServer(db database.DBLayer, storage storage.StorageLayer, tokenClient tokens.TokenClient, emiter msgqueue.EventEmiter) *Server {
	return &Server{
		DB:           db,
		Storage:      storage,
		MaxBodyBytes: MAX_BODY_BYTES,
		TokenClient:  tokenClient,
		Emitter:      emiter,
	}
}

// middleware for checking database connection
func (s *Server) CheckDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.DB == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "No database connection"})
			return
		}
		c.Next()
	}
}
