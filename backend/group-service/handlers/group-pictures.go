package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) SetGroupProfilePicture(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid ID"})
		return
	}
	groupID := c.Param("groupID")
	groupUID, err := uuid.Parse(groupID)
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
		c.JSON(http.StatusBadRequest, gin.H{"err": "image extention not allowed"})
		return
	}

	file, err := imageFileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad image"})
		return
	}

	pictureURL, err := s.DB.GetGroupProfilePictureURL(userUID, groupUID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
	}

	if err = s.Storage.UploadFile(file, pictureURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"newUrl": pictureURL})
}

func (s *Server) DeleteGroupProfilePicture(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
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

	pictureURL, err := s.DB.DeleteGroupProfilePicture(userUID, groupUID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	if err = s.Storage.DeleteFile(pictureURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
