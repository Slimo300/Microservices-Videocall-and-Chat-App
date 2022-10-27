package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/ws"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (s *Server) ServeWebSocket(c *gin.Context) {

	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	groups, err := s.DB.GetUserGroups(userUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ws.ServeWebSocket(c.Writer, c.Request, s.Hub, groups, userUID)
}

func (s *Server) AuthWS() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Query("authToken")
		if authToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"err": "no auth token provided"})
			return
		}
		accessToken, err := jwt.ParseWithClaims(authToken, &jwt.StandardClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return s.TokenService.GetPublicKey(), nil
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		userID := accessToken.Claims.(*jwt.StandardClaims).Subject
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Invalid token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
