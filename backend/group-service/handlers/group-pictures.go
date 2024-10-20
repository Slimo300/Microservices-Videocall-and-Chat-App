package handlers

import (
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app/command"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) SetGroupProfilePicture(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid ID"})
		return
	}
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid group ID"})
		return
	}

	imageFileHeader, err := c.FormFile("avatarFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	mimeType := imageFileHeader.Header.Get("Content-Type")
	if !isAllowedImageType(mimeType) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "image extension not allowed"})
		return
	}
	file, err := imageFileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad image"})
		return
	}
	if err := s.App.Commands.SetGroupPicture.Handle(c.Request.Context(), command.SetGroupPicture{UserID: userID, GroupID: groupID, File: file}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *Server) DeleteGroupProfilePicture(c *gin.Context) {
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

	if err := s.App.Commands.DeleteGroupPicture.Handle(c.Request.Context(), command.DeleteGroupPictureCommand{UserID: userID, GroupID: groupID}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
