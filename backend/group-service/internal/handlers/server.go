package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	"github.com/Slimo300/chat-groupservice/internal/database"
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
	"github.com/gin-gonic/gin"
)

type Server struct {
	DB           database.DBLayer
	Storage      storage.FileStorage
	TokenClient  tokens.TokenClient
	MaxBodyBytes int64
	Emitter      msgqueue.EventEmiter
	Listener     msgqueue.EventListener
}

func NewServer(db database.DBLayer, storage storage.FileStorage, tokenClient tokens.TokenClient) *Server {
	return &Server{
		DB:           db,
		Storage:      storage,
		MaxBodyBytes: 4194304,
		TokenClient:  tokenClient,
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
