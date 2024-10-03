package handlers

import (
	"fmt"
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app/command"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app/query"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getUserResponse struct {
	UserID     string `json:"ID"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	HasPicture bool   `json:"hasPicture"`
}

func (s *Server) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	user, err := s.app.Queries.GetUser.Handle(c.Request.Context(), query.GetUser{UserID: userID})
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	response := appUserToResponse(user)
	c.JSON(http.StatusOK, response)
}

type changePasswordRequest struct {
	OldPassword       string `json:"oldPassword"`
	NewPassword       string `json:"newPassword"`
	RepeatNewPassword string `json:"repeatPassword"`
}

func (s *Server) ChangePassword(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	var reqBody changePasswordRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if reqBody.NewPassword != reqBody.RepeatNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"err": "passwords don't match"})
		return
	}
	if err := s.app.Commands.ChangePassword.Handle(c.Request.Context(), command.ChangePassword{OldPassword: reqBody.OldPassword, NewPassword: reqBody.NewPassword, UserID: userID}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *Server) UpdateProfilePicture(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	imageFileHeader, err := c.FormFile("avatarFile")
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"err": fmt.Sprintf("Max request body size is %v bytes\n", s.maxBodyBytes),
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
	if err := s.app.Commands.SetProfilePicture.Handle(c.Request.Context(), command.SetProfilePicture{UserID: userID, File: file}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *Server) DeleteProfilePicture(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	if err := s.app.Commands.DeleteProfilePicture.Handle(c.Request.Context(), command.DeleteProfilePicture{UserID: userID}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func appUserToResponse(user models.User) getUserResponse {
	return getUserResponse{
		UserID:     user.ID().String(),
		Username:   user.Username(),
		Email:      user.Email(),
		HasPicture: user.HasPicture(),
	}
}

func isAllowedImageType(mimeType string) bool {
	var validImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	_, exists := validImageTypes[mimeType]
	return exists
}
