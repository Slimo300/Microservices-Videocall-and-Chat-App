package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/email"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
)

func (s *Server) RegisterUser(c *gin.Context) {
	load := struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Pass     string `json:"password"`
	}{}
	err := c.ShouldBindJSON(&load)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}
	if !isEmailValid(load.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid email"})
		return
	}
	if !isPasswordValid(load.Pass) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid password"})
		return
	}
	if len(load.UserName) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid username"})
		return
	}
	if s.DB.IsUsernameInDatabase(load.UserName) {
		c.JSON(http.StatusConflict, gin.H{"err": "username taken"})
		return
	}
	if s.DB.IsEmailInDatabase(load.Email) {
		c.JSON(http.StatusConflict, gin.H{"err": "email already in database"})
		return
	}
	load.Pass, err = hashPassword(load.Pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	user := models.User{Email: load.Email, UserName: load.UserName, Pass: load.Pass, Activated: false}
	user, err = s.DB.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	verificationCode, err := s.DB.NewVerificationCode(user.ID, randstr.String(10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if err := s.EmailService.SendVerificationEmail(email.VerificationEmailData{
		UserID:           user.ID.String(),
		Email:            user.Email,
		Name:             user.UserName,
		VerificationCode: verificationCode.ActivationCode,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

///////////////////////////////////////////////////////////////////////////////////////////////////
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
	requestPassword := load.Pass
	user, err := s.DB.GetUserByEmail(load.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "wrong email or password"})
		return
	}
	if !checkPassword(user.Pass, requestPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "wrong email or password"})
		return
	}
	if err := s.DB.SignInUser(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	tokenPair, err := s.AuthService.NewPairFromUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.SetCookie("refreshToken", tokenPair.RefreshToken, 86400, "/", s.Domain, false, true)

	c.JSON(http.StatusOK, gin.H{"accessToken": tokenPair.AccessToken})
}

///////////////////////////////////////////////////////////////////////////////////////////////
// SignOutUser method
func (s *Server) SignOutUser(c *gin.Context) {
	userID := c.GetString("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	if err := s.DB.SignOutUser(uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	refresh, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "No token to invalidate"})
		return
	}

	if err := s.AuthService.DeleteUserToken(refresh); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.SetCookie("refreshToken", "", -1, "/", s.Domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

///////////////////////////////////////////////////////////////////////////////////////////////
// GetUserById method
func (s *Server) GetUser(c *gin.Context) {
	userID := c.GetString("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
	}

	user, err := s.DB.GetUserById(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no such user"})
		return
	}
	user.Pass = ""

	c.JSON(http.StatusOK, user)
}

func (s *Server) RefreshToken(c *gin.Context) {

	refresh, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "No token provided"})
		return
	}
	tokens, err := s.AuthService.NewPairFromRefresh(refresh)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if tokens.Error != "" {
		if tokens.Error == "Token Blacklisted" {
			c.JSON(http.StatusForbidden, gin.H{"err": "Token Blacklisted"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": tokens.Error})
		return
	}
	c.SetCookie("refreshToken", tokens.RefreshToken, 86400, "/", s.Domain, false, true)

	c.JSON(http.StatusOK, gin.H{"accessToken": tokens.AccessToken})
}
