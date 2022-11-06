package handlers

import (
	"log"
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/events"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetUserGroups(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
	}

	groups, err := s.DB.GetUserGroups(userUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(groups) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, groups)

}

func (s *Server) CreateGroup(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
	}

	var group models.Group
	err = c.ShouldBindJSON(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if group.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad name"})
		return
	}
	if group.Desc == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad description"})
		return
	}

	group, err = s.DB.CreateGroup(userUID, group.Name, group.Desc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if err := s.Emitter.Emit(events.MemberCreatedEvent{
		ID:      group.Members[0].ID,
		GroupID: group.ID,
		UserID:  userUID,
		Creator: true,
	}); err != nil {
		log.Printf("Emitter failed: %s", err.Error())
	}

	c.JSON(http.StatusCreated, group)
}

func (s *Server) DeleteGroup(c *gin.Context) {
	userID := c.GetString("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	groupID := c.Param("groupID")
	groupUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}

	member, err := s.DB.GetUserGroupMember(uid, groupUID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "couldn't delete group"})
		return
	}
	if !member.Creator {
		c.JSON(http.StatusForbidden, gin.H{"err": "couldn't delete group"})
		return
	}

	group, err := s.DB.DeleteGroup(groupUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	if err := s.Emitter.Emit(events.GroupDeletedEvent{
		ID: group.ID,
	}); err != nil {
		log.Printf("Emitter failed: %s", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}
