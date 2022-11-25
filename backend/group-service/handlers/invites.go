package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
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

func (s *Server) CreateInvite(c *gin.Context) {
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
	targetUUID, err := uuid.Parse(load.Target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid target user ID"})
		return
	}

	invite, err := s.DB.AddInvite(userUID, targetUUID, groupUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal database error"})
		return
	}

	s.Emitter.Emit(events.InviteSentEvent{
		ID:       invite.ID,
		IssuerID: invite.IssId,
		TargetID: invite.TargetID,
	})

	c.JSON(http.StatusCreated, gin.H{"message": "invite sent"})
}

func (s *Server) RespondGroupInvite(c *gin.Context) {
	userID := c.GetString("userID")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	inviteID := c.Param("inviteID")
	inviteUUID, err := uuid.Parse(inviteID)
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

	invite, group, member, err := s.DB.AnswerInvite(userUUID, inviteUUID, *load.Answer)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error})
		return
	}

	if member != nil {
		s.Emitter.Emit(events.MemberCreatedEvent{
			ID:      member.ID,
			GroupID: member.GroupID,
			UserID:  member.UserID,
			Creator: false,
		})
	}
	if invite != nil {
		var answer bool
		if invite.Status == models.INVITE_ACCEPT {
			answer = true
		}
		// Emit InviteUpdate
		s.Emitter.Emit(events.InviteRespondedEvent{
			ID:       invite.ID,
			IssuerID: invite.IssId,
			Answer:   answer,
		})
	}

	c.JSON(http.StatusOK, group)
}
