package handlers

import (
	"net/http"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/ws"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
)

func (s *Server) ServeWebSocket(c *gin.Context) {

	accessCode := c.Query("accessCode")
	if accessCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no access code provided"})
		return
	}

	user, err := s.CodeCache.Read(accessCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "access code not found"})
		return
	}
	if user.Deadline.Unix() < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"err": "access code invalidated"})
		return
	}
	s.CodeCache.Delete(accessCode)

	groups, err := s.DB.GetUserGroups(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ws.ServeWebSocket(c.Writer, c.Request, *s.Hub, groups, user.ID)
}

func (s *Server) GetAuthCode(c *gin.Context) {

	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	newCode := randstr.String(10)
	s.CodeCache.Set(newCode, userUID)

	c.JSON(http.StatusOK, gin.H{"accessCode": newCode})
}

// func (s *Server) AuthWS() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authToken := c.Query("authToken")
// 		if authToken == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"err": "no auth token provided"})
// 			return
// 		}
// 		accessToken, err := jwt.ParseWithClaims(authToken, &jwt.StandardClaims{},
// 			func(t *jwt.Token) (interface{}, error) {
// 				return s.TokenService.GetPublicKey(), nil
// 			})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
// 			return
// 		}
// 		userID := accessToken.Claims.(*jwt.StandardClaims).Subject
// 		if userID == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"err": "Invalid token"})
// 			return
// 		}

// 		c.Set("userID", userID)
// 		c.Next()
// 	}
// }
