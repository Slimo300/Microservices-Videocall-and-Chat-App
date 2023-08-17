package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/email"
	"github.com/gin-gonic/gin"
)

func (s *Server) ForgotPassword(c *gin.Context) {
	queryEmail := c.Query("email")
	if !isEmailValid(queryEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid email address"})
		return
	}

	user, resetCode, err := s.DB.NewResetPasswordCode(queryEmail)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	if resetCode == nil {
		c.JSON(http.StatusOK, gin.H{"message": "reset password email sent"})
	}

	if user != nil && resetCode != nil {
		go func() {
			if _, err := s.EmailClient.SendResetPasswordEmail(context.TODO(), &email.EmailData{
				Email: user.Email,
				Name:  user.UserName,
				Code:  resetCode.ResetCode,
			}); err != nil {
				log.Println(err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "reset password email sent"})
}

func (s *Server) ResetForgottenPassword(c *gin.Context) {

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusForbidden, "invalid reset code")
		return
	}

	payload := struct {
		NewPassword    string `json:"newPassword"`
		RepeatPassword string `json:"repeatPassword"`
	}{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if !isPasswordValid(payload.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid password"})
		return
	}
	if payload.NewPassword != payload.RepeatPassword {
		c.JSON(http.StatusBadRequest, gin.H{"err": "passwords don't match"})
		return
	}

	if err := s.DB.ResetPassword(code, payload.NewPassword); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}
