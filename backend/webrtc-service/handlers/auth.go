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
	member, err := s.DB.GetMemberByGroupAndUserID(groupID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Generate auth code and insert it to database
	accessCode := randstr.String(10)
	if err := s.DB.NewAccessCode(accessCode, member.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	// Return auth code
	c.JSON(http.StatusOK, gin.H{
		"accessCode": accessCode,
		"username":   member.Username,
		"muting":     member.Muting,
	})
}
