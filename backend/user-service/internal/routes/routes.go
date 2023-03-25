package routes

import (
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
	"github.com/Slimo300/chat-userservice/internal/handlers"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func Setup(server *handlers.Server, origin string) *gin.Engine {
	engine := gin.Default()

	engine.Use(CORSMiddleware(origin))
	engine.Use(limits.RequestSizeLimiter(server.MaxBodyBytes))

	api := engine.Group("/users")
	api.POST("/register", server.RegisterUser)
	api.GET("/verify-account/:code", server.VerifyCode)

	api.GET("/forgot-password", server.ForgotPassword)
	api.PATCH("/reset-password/:code", server.ResetForgottenPassword)

	api.POST("/login", server.SignIn)
	api.POST("/refresh", server.RefreshToken)

	apiAuth := api.Use(tokens.MustAuth(server.TokenClient))

	apiAuth.POST("/signout", server.SignOutUser)

	apiAuth.DELETE("/delete-image", server.DeleteProfilePicture)
	apiAuth.POST("/set-image", server.UpdateProfilePicture)

	apiAuth.PUT("/change-password", server.ChangePassword)
	apiAuth.GET("/user", server.GetUser)

	return engine
}

func CORSMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
