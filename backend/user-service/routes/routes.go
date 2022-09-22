package routes

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/handlers"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, server *handlers.Server) {

	engine.Use(CORSMiddleware())
	engine.Use(limits.RequestSizeLimiter(server.MaxBodyBytes))

	api := engine.Group("/api")
	api.POST("/register", server.RegisterUser)
	api.POST("/login", server.SignIn)
	api.POST("/refresh", server.RefreshToken)

	apiAuth := api.Use(auth.MustAuth(server.AuthService))

	apiAuth.DELETE("/delete-image", server.DeleteProfilePicture)
	apiAuth.POST("/set-image", server.UpdateProfilePicture)
	apiAuth.PUT("/change-password", server.ChangePassword)
	apiAuth.POST("/signout", server.SignOutUser)
	apiAuth.GET("/user", server.GetUser)

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
