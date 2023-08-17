package auth

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func MustAuthWithKey(key *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessHeader := strings.Split(c.GetHeader("Authorization"), " ")[1]
		if accessHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "user not authenticated"})
			return
		}
		accessToken, err := jwt.ParseWithClaims(accessHeader, &jwt.StandardClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return key, nil
			})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Invalid token"})
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
