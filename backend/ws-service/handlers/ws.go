package handlers

import (
	"net/http"

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

	userID, err := s.DB.CheckAccessCode(accessCode)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "connection not authorized"})
		return
	}

	groups, err := s.DB.GetUserGroups(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ws.ServeWebSocket(c.Writer, c.Request, s.Hub, groups, userID)
}

func (s *Server) GetAuthCode(c *gin.Context) {

	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	newCode := randstr.String(10)
	if err := s.DB.NewAccessCode(userUID, newCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessCode": newCode})
}
