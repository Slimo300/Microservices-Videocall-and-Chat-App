package handlers

import (
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/google/uuid"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app/command"
	"github.com/gin-gonic/gin"
)

type registerUserRequest struct {
	UserName       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"rpassword"`
}

func (s *Server) RegisterUser(c *gin.Context) {
	var reqBody registerUserRequest
	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if reqBody.Password != reqBody.RepeatPassword {
		c.JSON(http.StatusBadRequest, gin.H{"err": "passwords don't match"})
		return
	}
	if err := s.app.Commands.RegisterUser.Handle(c.Request.Context(), command.RegisterUser{
		Email:    reqBody.Email,
		Username: reqBody.UserName,
		Password: reqBody.Password,
	}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (s *Server) VerifyCode(c *gin.Context) {
	code, err := uuid.Parse(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "code not found"})
		return
	}
	if err := s.app.Commands.VerifyEmail.Handle(c.Request.Context(), command.VerifyEmail{Code: code}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
