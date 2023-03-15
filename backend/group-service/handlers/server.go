package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	DB           database.DBLayer
	Storage      storage.FileStorage
	TokenService auth.TokenClient
	MaxBodyBytes int64
	Emitter      msgqueue.EventEmiter
	Listener     msgqueue.EventListener
}

func NewServer(db database.DBLayer, storage storage.FileStorage, auth auth.TokenClient) *Server {
	return &Server{
		DB:           db,
		Storage:      storage,
		MaxBodyBytes: 4194304,
		TokenService: auth,
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
