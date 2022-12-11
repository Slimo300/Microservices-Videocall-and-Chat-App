package handlers

import (
	"fmt"
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////
// ChangePassword
func (s *Server) ChangePassword(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	payload := struct {
		NewPassword       string `json:"newPassword"`
		RepeatNewPassword string `json:"repeatPassword"`
		OldPassword       string `json:"oldPassword"`
	}{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if !isPasswordValid(payload.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Password must be at least 6 characters long"})
		return
	}
	if payload.NewPassword != payload.RepeatNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Passwords don't match"})
		return
	}

	if err := s.DB.ChangePassword(userUID, payload.OldPassword, payload.NewPassword); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

///////////////////////////////////////////////////////////////////////////
// UpdateProfilePicture

func (s *Server) UpdateProfilePicture(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	imageFileHeader, err := c.FormFile("avatarFile")
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"err": fmt.Sprintf("Max request body size is %v bytes\n", s.MaxBodyBytes),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if imageFileHeader == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no file provided"})
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

	pictureURL, err := s.DB.GetProfilePictureURL(userUID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	if err = s.ImageStorage.UpdateProfilePicture(file, pictureURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"newUrl": pictureURL})
}

///////////////////////////////////////////////////////////////////////////
// UpdateProfilePicture

func (s *Server) DeleteProfilePicture(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	url, err := s.DB.DeleteProfilePicture(userUID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": "User not found"})
		return
	}
	if err = s.ImageStorage.DeleteProfilePicture(url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
