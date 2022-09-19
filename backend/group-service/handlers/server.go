package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/auth"
	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/communication"
	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/database"
	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/storage"
	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/ws"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Server struct {
	DB           database.DBlayer
	Storage      storage.StorageLayer
	Hub          ws.HubInterface
	TokenService auth.TokenClient
	actionChan   chan<- *communication.Action
	messageChan  <-chan *communication.Message
	domain       string
	MaxBodyBytes int64
}

func NewServer(db database.DBlayer, storage storage.StorageLayer, auth auth.TokenClient) *Server {
	actionChan := make(chan *communication.Action)
	messageChan := make(chan *communication.Message)
	return &Server{
		DB:           db,
		Storage:      storage,
		domain:       "localhost",
		actionChan:   actionChan,
		messageChan:  messageChan,
		MaxBodyBytes: 4194304,
		TokenService: auth,
		Hub:          ws.NewHub(messageChan, actionChan),
	}
}

func NewServerWithMockHub(db database.DBlayer, storage storage.StorageLayer) *Server {
	actionChan := make(chan *communication.Action)
	messageChan := make(chan *communication.Message)
	return &Server{
		DB:           db,
		Storage:      storage,
		domain:       "localhost",
		actionChan:   actionChan,
		messageChan:  messageChan,
		MaxBodyBytes: 4194304,
		Hub:          ws.NewMockHub(actionChan),
	}
}

func (s *Server) RunHub() {
	go s.ListenToHub()
	s.Hub.Run()
}

func (s *Server) ListenToHub() {
	var msg *communication.Message
	for {
		select {
		case msg = <-s.messageChan:
			when, err := time.Parse(communication.TIME_FORMAT, msg.When)
			if err != nil {
				panic(err.Error())
			}
			if err := s.DB.AddMessage(msg.Member, msg.Message, when); err != nil {
				panic("Panicked while adding message")
			}
		}
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
