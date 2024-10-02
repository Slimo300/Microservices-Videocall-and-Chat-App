package handlers

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func (s *Server) setup(origin string) *gin.Engine {
	engine := gin.Default()

	engine.Use(corsMiddleware(origin))
	engine.Use(limits.RequestSizeLimiter(s.maxBodyBytes))
	engine.Use(func(ctx *gin.Context) {
		if len(ctx.Errors) > 0 {
			return
		}
	})

	api := engine.Group("/users")
	api.POST("/register", s.RegisterUser)
	api.GET("/verify-account/:code", s.VerifyCode)

	api.GET("/forgot-password", s.ForgotPassword)
	api.PATCH("/reset-password/:code", s.ResetForgottenPassword)

	api.POST("/login", s.SignIn)
	api.POST("/refresh", s.RefreshToken)

	apiAuth := api.Use(auth.MustAuthWithKey(s.publicKey))

	apiAuth.POST("/signout", s.SignOutUser)

	apiAuth.DELETE("/delete-image", s.DeleteProfilePicture)
	apiAuth.POST("/set-image", s.UpdateProfilePicture)

	apiAuth.PUT("/change-password", s.ChangePassword)
	apiAuth.GET("/user", s.GetUser)

	return engine
}

func corsMiddleware(origin string) gin.HandlerFunc {
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
