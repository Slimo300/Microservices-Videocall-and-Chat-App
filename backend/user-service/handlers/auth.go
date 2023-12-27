package handlers

import (
	"context"
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	"github.com/gin-gonic/gin"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////
// SignIn method
func (s *Server) SignIn(c *gin.Context) {
	load := struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}{}
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	user, err := s.DB.SignIn(load.Email, load.Pass)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	tokenPair, err := s.TokenClient.NewPairFromUserID(context.TODO(), &auth.UserID{ID: user.ID.String()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if tokenPair.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": tokenPair.Error})
		return
	}

	c.SetCookie("refreshToken", tokenPair.RefreshToken, 86400, "/", s.Domain, true, true)

	c.JSON(http.StatusOK, gin.H{"accessToken": tokenPair.AccessToken})
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// SignOutUser method
func (s *Server) SignOutUser(c *gin.Context) {

	refresh, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "No token to invalidate"})
		return
	}

	if _, err := s.TokenClient.DeleteUserToken(context.TODO(), &auth.RefreshToken{Token: refresh}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.SetCookie("refreshToken", "", -1, "/", s.Domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// //////////////////////////////////////////////////////////////////////////////////////////////
// Refresh Token
func (s *Server) RefreshToken(c *gin.Context) {

	refresh, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "No token provided"})
		return
	}
	tokens, err := s.TokenClient.NewPairFromRefresh(context.TODO(), &auth.RefreshToken{Token: refresh})
	if tokens.Error != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": tokens.Error})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	// Deletes refreshToken cookie. If token verification won't fail cookie will be repopulated
	c.SetCookie("refreshToken", "", -1, "/", s.Domain, false, true)

	if tokens.Error != "" {
		if tokens.Error == "Token Blacklisted" {
			c.JSON(http.StatusForbidden, gin.H{"err": tokens.Error})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": tokens.Error})
		return
	}
	c.SetCookie("refreshToken", tokens.RefreshToken, 86400, "/", s.Domain, false, true)

	c.JSON(http.StatusOK, gin.H{"accessToken": tokens.AccessToken})
}
