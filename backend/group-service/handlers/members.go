package handlers

import (
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GrantRights(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid member ID"})
		return
	}

	var rights models.MemberRights
	if err := c.ShouldBindJSON(&rights); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if rights.Adding == 0 && rights.DeletingMessages == 0 && rights.DeletingMembers == 0 && rights.Admin == 0 && rights.Muting == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no action specified"})
		return
	}

	_, err = s.Service.GrantRights(userID, memberID, rights)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member updated"})
}

func (s *Server) DeleteUserFromGroup(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid member ID"})
		return
	}

	_, err = s.Service.DeleteMember(userID, memberID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member deleted"})
}
