package handlers

import (
	"log"
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/email"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterUser(c *gin.Context) {
	payload := struct {
		UserName   string `json:"username"`
		Email      string `json:"email"`
		Pass       string `json:"password"`
		RepeatPass string `json:"rpassword"`
	}{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}
	if !isEmailValid(payload.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid email"})
		return
	}
	if len(payload.UserName) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid username"})
		return
	}
	if !isPasswordValid(payload.Pass) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid password"})
		return
	}
	if payload.Pass != payload.RepeatPass {
		c.JSON(http.StatusBadRequest, gin.H{"err": "passwords don't match"})
		return
	}
	user, verificationCode, err := s.DB.RegisterUser(models.User{Email: payload.Email, UserName: payload.UserName, Pass: payload.Pass, Verified: false})
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	if user != nil && verificationCode != nil {
		go func() {
			s.EmailService.SendVerificationEmail(email.EmailData{
				Subject: "Verification Email",
				Email:   user.Email,
				Name:    user.UserName,
				Code:    verificationCode.ActivationCode,
				Origin:  s.Origin,
			})
		}()
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (s *Server) VerifyCode(c *gin.Context) {
	code := c.Param("code")

	user, err := s.DB.VerifyCode(code)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	if err := s.Emitter.Emit(events.UserRegisteredEvent{
		ID:         user.ID,
		Username:   user.UserName,
		PictureURL: user.PictureURL,
	}); err != nil {
		log.Printf("Emitter error: %v", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"message": "code verified"})
}
