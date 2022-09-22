package auth

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/pb"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenClient interface {
	NewPairFromUserID(userID uuid.UUID) (*pb.TokenPair, error)
	NewPairFromRefresh(refresh string) (*pb.TokenPair, error)
	DeleteUserToken(refresh string) error
	GetPublicKey() *rsa.PublicKey
}

func MustAuth(auth TokenClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessHeader := strings.Split(c.GetHeader("Authorization"), " ")[1]
		if accessHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "user not authenticated"})
			return
		}
		accessToken, err := jwt.ParseWithClaims(accessHeader, &jwt.StandardClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return auth.GetPublicKey(), nil
			})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		userID := accessToken.Claims.(*jwt.StandardClaims).Subject
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Invalid token"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
