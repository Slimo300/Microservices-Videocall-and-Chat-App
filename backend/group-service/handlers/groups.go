package handlers

import (
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app/command"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app/query"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetUserGroups(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
	}
	groups, err := s.App.Queries.GetUserGroups.Handle(c.Request.Context(), query.GetUserGroups{UserID: userID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

type createGroupRequest struct {
	Name string `json:"name"`
}

func (s *Server) CreateGroup(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	var reqBody createGroupRequest
	err = c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if reqBody.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "name not specified"})
		return
	}
	if err := s.App.Commands.CreateGroup.Handle(c.Request.Context(), command.CreateGroupCommand{UserID: userID, Name: reqBody.Name}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (s *Server) DeleteGroup(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}
	if err = s.App.Commands.DeleteGroup.Handle(c.Request.Context(), command.DeleteGroupCommand{UserID: userID, GroupID: groupID}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
}
