package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/app/query"
)

type presignedGetRequestBody struct {
	Files []storage.PresignGetFileInput `json:"files"`
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
	presignedRequests, err := s.App.Queries.GetPresignedGetRequestsHandler.Handle(c.Request.Context(), query.GetPresignedGetRequestsQuery{
		UserID:   userID,
		GroupID:  groupID,
		FileKeys: files.Files,
	})
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, presignedRequests)
}

type presignedPutRequestBody struct {
	Files []storage.PresignPutFileInput `json:"files"`
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
	presignedRequests, err := s.App.Queries.GetPresignedPutRequestsHandler.Handle(c.Request.Context(), query.GetPresignedPutRequestsQuery{
		UserID:   userID,
		GroupID:  groupID,
		FileKeys: files.Files,
	})
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, presignedRequests)
}
