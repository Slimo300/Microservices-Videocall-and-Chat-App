package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/storage"
)

type presignedGetRequestBody struct {
	Files []storage.GetFileInput `json:"files"`
}

func (s *Server) GetPresignedGetRequests(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid user ID"})
		return
	}
	groupID, err := uuid.Parse((c.Param("groupID")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}

	var files presignedGetRequestBody
	if err := c.ShouldBindJSON(&files); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(files.Files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid request body"})
		return
	}

	if _, err := s.DB.GetGroupMembership(userID, groupID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "user cannot send messages to this group"})
		return
	}

	presignedRequests, err := s.Storage.GetPresignedGetRequests(groupID.String(), files.Files...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, presignedRequests)
}

type presignedPutRequestBody struct {
	Files []storage.PutFileInput `json:"files"`
}

func (s *Server) GetPresignedPutRequests(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid user ID"})
		return
	}
	groupID, err := uuid.Parse((c.Param("groupID")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}

	var files presignedPutRequestBody
	if err := c.ShouldBindJSON(&files); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(files.Files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid request body"})
		return
	}

	if _, err := s.DB.GetGroupMembership(userID, groupID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "user cannot send messages to this group"})
		return
	}

	requestsData, err := s.Storage.GetPresignedPutRequests(groupID.String(), files.Files...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requestsData)

}
