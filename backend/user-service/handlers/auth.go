package handlers

import (
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app/command"
	"github.com/gin-gonic/gin"
)

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) SignIn(c *gin.Context) {
	var reqBody signInRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	accessToken, refreshToken, err := s.app.Commands.SignIn.Handle(c.Request.Context(), command.SignIn{Email: reqBody.Email, Password: reqBody.Password})
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.SetCookie("refreshToken", refreshToken, 86400, "/", s.domain, true, true)
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

func (s *Server) SignOutUser(c *gin.Context) {
	refresh, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "No token to invalidate"})
		return
	}
	if err := s.app.Commands.SignOut.Handle(c.Request.Context(), command.SignOut{RefreshToken: refresh}); err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.SetCookie("refreshToken", "", -1, "/", s.domain, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s *Server) RefreshToken(c *gin.Context) {
	refresh, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "No token provided"})
		return
	}
	accessToken, refreshToken, err := s.app.Commands.RefreshToken.Handle(c.Request.Context(), command.RefreshToken{RefreshToken: refresh})
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	c.SetCookie("refreshToken", refreshToken, 86400, "/", s.domain, false, true)
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}
