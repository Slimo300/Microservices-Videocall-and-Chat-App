package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/communication"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GrantPriv(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	memberID := c.Param("memberID")
	memberUID, err := uuid.Parse(memberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid member ID"})
		return
	}

	load := struct {
		Adding           *bool `json:"adding" binding:"required"`
		DeletingMessages *bool `json:"deleting" binding:"required"`
		DeletingMembers  *bool `json:"deletingMembers" binding:"required"`
		Setting          *bool `json:"setting" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad request, all 3 fields must be present"})
		return
	}

	memberToBeChanged, err := s.DB.GetMemberByID(memberUID)
	if err != nil || memberToBeChanged.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"err": "resource not found"})
		return
	}

	issuerMember, err := s.DB.GetUserGroupMember(userUID, memberToBeChanged.GroupID)
	if err != nil || !issuerMember.Setting {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to put"})
		return
	}

	if memberToBeChanged.Creator {
		c.JSON(http.StatusForbidden, gin.H{"err": "creator can't be modified"})
	}

	if err := s.DB.GrantPriv(memberUID, *load.Adding, *load.DeletingMembers, *load.Setting, *load.DeletingMessages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (s *Server) DeleteUserFromGroup(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	memberID := c.Param("memberID")
	memberUID, err := uuid.Parse(memberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid member ID"})
		return
	}

	memberToBeDeleted, err := s.DB.GetMemberByID(memberUID)
	if err != nil || memberToBeDeleted.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"err": "resource not found"})
		return
	}

	issuerMember, err := s.DB.GetUserGroupMember(userUID, memberToBeDeleted.GroupID)
	if err != nil || (!issuerMember.DeletingMembers && issuerMember.ID != memberToBeDeleted.ID) {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to delete"})
		return
	}

	member, err := s.DB.DeleteUserFromGroup(memberUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.actionChan <- &communication.Action{Action: "DELETE_MEMBER", Member: member}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
