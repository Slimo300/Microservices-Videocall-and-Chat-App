package handlers

import (
	"net/http"
	"strconv"

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
	num, err := strconv.Atoi(c.Query("num"))
	if err != nil || num <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "number of messages is not a valid number"})
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "offset is not a valid number"})
		return
	}
	invites, err := s.DB.GetUserInvites(userUID, num, offset)
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

	payload := struct {
		GroupID string `json:"group"`
		Target  string `json:"target"`
	}{}

	// getting req body
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	groupUID, err := uuid.Parse(payload.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}
	targetUUID, err := uuid.Parse(payload.Target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid target user ID"})
		return
	}

	invite, err := s.DB.AddInvite(userUID, targetUUID, groupUID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	_ = s.Emitter.Emit(events.InviteSentEvent{
		ID:       invite.ID,
		IssuerID: invite.IssId,
		Issuer: events.User{
			UserName: invite.Iss.UserName,
			Picture:  invite.Iss.Picture,
		},
		TargetID: invite.TargetID,
		GroupID:  invite.GroupID,
		Group: events.Group{
			Name:    invite.Group.Name,
			Picture: invite.Group.Picture,
		},
		Status:   int(invite.Status),
		Modified: invite.Modified,
	})

	c.JSON(http.StatusCreated, invite)
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

	payload := struct {
		Answer *bool `json:"answer" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "answer not specified"})
		return
	}

	invite, group, member, err := s.DB.AnswerInvite(userUUID, inviteUUID, *payload.Answer)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	if member != nil {
		_ = s.Emitter.Emit(events.MemberCreatedEvent{
			ID:      member.ID,
			GroupID: member.GroupID,
			UserID:  member.UserID,
			User: events.User{
				UserName: member.User.UserName,
				Picture:  member.User.Picture,
			},
			Adding:           member.Adding,
			DeletingMembers:  member.DeletingMembers,
			DeletingMessages: member.DeletingMessages,
			Admin:            member.Admin,
			Creator:          member.Creator,
		})
	}
	if invite != nil {
		_ = s.Emitter.Emit(events.InviteRespondedEvent{
			ID:       invite.ID,
			IssuerID: invite.IssId,
			TargetID: invite.TargetID,
			Target: events.User{
				UserName: invite.Target.UserName,
				Picture:  invite.Target.Picture,
			},
			GroupID: invite.GroupID,
			Group: events.Group{
				Name:    invite.Group.Name,
				Picture: invite.Group.Picture,
			},
			Status:   int(invite.Status),
			Modified: invite.Modified,
		})
	}
	if !*payload.Answer {
		c.JSON(http.StatusOK, gin.H{"invite": invite})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invite": invite, "group": group})
}
