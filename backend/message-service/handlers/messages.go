package handlers

import (
	"net/http"
	"strconv"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/app/command"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/app/query"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetGroupMessages(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
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
	messages, err := s.App.Queries.GetGroupMessages.Handle(c.Request.Context(), query.GetGroupMessagesQuery{
		UserID:  userID,
		GroupID: groupID,
		Num:     num,
		Offset:  offset,
	})
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (s *Server) DeleteMessageForEveryone(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid user ID"})
		return
	}
	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid message ID"})
		return
	}
	if err := s.App.Commands.DeleteMessageForEveryone.Handle(c.Request.Context(), command.DeleteMessageForEveryoneCommand{
		MessageID: messageID,
		UserID:    userID,
	}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *Server) DeleteMessageForYourself(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid user ID"})
		return
	}
	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid message ID"})
		return
	}
	if err := s.App.Commands.DeleteMessageForYourself.Handle(c.Request.Context(), command.DeleteMessageForYourselfCommand{
		MessageID: messageID,
		UserID:    userID,
	}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
