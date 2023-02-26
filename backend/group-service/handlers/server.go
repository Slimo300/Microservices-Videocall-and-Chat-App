package handlers

import (
	"net/http"
	"strings"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

func (s *Server) MustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessHeader := strings.Split(c.GetHeader("Authorization"), " ")[1]
		if accessHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "user not authenticated"})
			return
		}
		accessToken, err := jwt.ParseWithClaims(accessHeader, &jwt.StandardClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return s.TokenService.GetPublicKey(), nil
			})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		userID := accessToken.Claims.(*jwt.StandardClaims).Subject
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Invalid token"})
			return
		}
		c.Set("userID", userID)
		c.Next()
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
