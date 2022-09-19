package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetUserInvites(c *gin.Context) {

	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid ID"})
		return
	}

	invites, err := s.DB.GetUserInvites(userUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if len(invites) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, invites)
}

func (s *Server) SendGroupInvite(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	load := struct {
		GroupID string `json:"group"`
		Target  string `json:"target"`
	}{}

	// getting req body
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	groupUID, err := uuid.Parse(load.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}
	if strings.TrimSpace(load.Target) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not specified"})
		return
	}

	issuerMember, err := s.DB.GetUserGroupMember(userUID, groupUID)
	if err != nil || !issuerMember.Adding {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to add"})
		return
	}

	userToBeAdded, err := s.DB.GetUserByUsername(load.Target)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": fmt.Sprintf("no user with name: %s", load.Target)})
		return
	}

	if s.DB.IsUserInGroup(userToBeAdded.ID, groupUID) {
		c.JSON(http.StatusConflict, gin.H{"err": "user is already a member of group"})
		return
	}

	if s.DB.IsUserInvited(userToBeAdded.ID, groupUID) {
		c.JSON(http.StatusConflict, gin.H{"err": "user already invited"})
		return
	}

	_, err = s.DB.AddInvite(userUID, userToBeAdded.ID, groupUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal database error"})
		return
	}

	// s.actionChan <- &communication.Action{Invite: invite}

	c.JSON(http.StatusCreated, gin.H{"message": "invite sent"})
}

func (s *Server) RespondGroupInvite(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	inviteID := c.Param("inviteID")
	inviteUID, err := uuid.Parse(inviteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid invite id"})
		return
	}

	load := struct {
		Answer *bool `json:"answer" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "answer not specified"})
		return
	}

	invite, err := s.DB.GetInviteByID(inviteUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "resource not found"})
		return
	}

	if invite.TargetID != userUID {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to respond"})
		return
	}

	if invite.Status != database.INVITE_AWAITING {
		c.JSON(http.StatusForbidden, gin.H{"err": "invite already answered"})
		return
	}

	if !*load.Answer {
		if err := s.DB.DeclineInvite(invite.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "invite declined"})
		return
	}

	group, err := s.DB.AcceptInvite(invite)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no such invite"})
		return
	}

	c.JSON(http.StatusOK, group)
}
