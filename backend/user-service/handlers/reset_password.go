package handlers

import (
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app/command"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) ForgotPassword(c *gin.Context) {
	queryEmail := c.Query("email")
	if err := s.app.Commands.ForgotPassword.Handle(c.Request.Context(), command.ForgotPassword{Email: queryEmail}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

type resetForgottenPasswordRequest struct {
	NewPassword    string `json:"newPassword"`
	RepeatPassword string `json:"repeatPassword"`
}

func (s *Server) ResetForgottenPassword(c *gin.Context) {
	code, err := uuid.Parse(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "code not found"})
		return
	}
	var reqBody resetForgottenPasswordRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if reqBody.NewPassword != reqBody.RepeatPassword {
		c.JSON(http.StatusBadRequest, gin.H{"err": "passwords don't match"})
		return
	}
	if err := s.app.Commands.ResetForgottenPassword.Handle(c.Request.Context(), command.ResetForgottenPassword{Code: code, NewPassword: reqBody.NewPassword}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
