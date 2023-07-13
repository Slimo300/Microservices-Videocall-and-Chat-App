package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
)

func (s *Server) GetAuthCode(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid user ID"})
		return
	}

	groupID := c.Param("groupID")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "groupID not provided"})
		return
	}

	// Check if user belongs to a group
	if _, err := s.DB.GetMember(userID, groupID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}

	// Check if group has a session started

	// Generate auth code and insert it to database
	accessCode := randstr.String(10)
	if err := s.DB.NewAccessCode(groupID, userID, accessCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	// Return auth code
	c.JSON(http.StatusOK, gin.H{"accessCode": accessCode})
}
