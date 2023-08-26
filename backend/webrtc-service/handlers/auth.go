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

	memberID := c.Param("memberID")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "groupID not provided"})
		return
	}

	// Check if user belongs to a group
	member, err := s.DB.GetMember(memberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid member ID"})
		return
	}

	if member.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"err": "Membership does not belong to authenticated user"})
		return
	}

	// Generate auth code and insert it to database
	accessCode := randstr.String(10)
	if err := s.DB.NewAccessCode(accessCode, memberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	// Return auth code
	c.JSON(http.StatusOK, gin.H{"accessCode": accessCode})
}
