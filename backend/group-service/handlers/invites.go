package handlers

import (
	"net/http"
	"strconv"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetUserInvites(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
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

	invites, err := s.Service.GetUserInvites(c.Request.Context(), userID, num, offset)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	if len(invites) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, invites)
}

func (s *Server) CreateInvite(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
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

	groupID, err := uuid.Parse(payload.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}

	targetID, err := uuid.Parse(payload.Target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid target user ID"})
		return
	}

	invite, err := s.Service.AddInvite(c.Request.Context(), userID, targetID, groupID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invite)
}

func (s *Server) RespondGroupInvite(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	inviteID, err := uuid.Parse(c.Param("inviteID"))
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

	invite, group, err := s.Service.RespondInvite(c.Request.Context(), userID, inviteID, *payload.Answer)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	if !*payload.Answer {
		c.JSON(http.StatusOK, gin.H{"invite": invite})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invite": invite, "group": group})
}
