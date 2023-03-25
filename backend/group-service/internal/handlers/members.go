package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/chat-groupservice/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GrantPriv(c *gin.Context) {
	userID := c.GetString("userID")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	groupID := c.Param("groupID")
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}
	memberID := c.Param("memberID")
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid member ID"})
		return
	}

	var rights models.MemberRights
	if err := c.ShouldBindJSON(&rights); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if rights.Adding == 0 && rights.DeletingMessages == 0 && rights.DeletingMembers == 0 && rights.Admin == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no action specified"})
		return
	}

	member, err := s.DB.GrantRights(userUUID, groupUUID, memberUUID, rights)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	if member != nil {
		_ = s.Emitter.Emit(events.MemberUpdatedEvent{
			ID:      member.ID,
			GroupID: member.GroupID,
			UserID:  member.UserID,
			User: events.User{
				UserName: member.User.UserName,
				Picture:  member.User.Picture,
			},
			DeletingMessages: member.DeletingMessages,
			DeletingMembers:  member.DeletingMembers,
			Adding:           member.Adding,
			Admin:            member.Admin,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "member updated"})
}

func (s *Server) DeleteUserFromGroup(c *gin.Context) {
	userID := c.GetString("userID")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	groupID := c.Param("groupID")
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}
	memberID := c.Param("memberID")
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid member ID"})
		return
	}

	member, err := s.DB.DeleteMember(userUUID, groupUUID, memberUUID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	_ = s.Emitter.Emit(events.MemberDeletedEvent{ID: member.ID, GroupID: member.GroupID, UserID: member.UserID})

	c.JSON(http.StatusOK, gin.H{"message": "member deleted"})
}
