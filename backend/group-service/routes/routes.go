package routes

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func Setup(server *handlers.Server, origin string) *gin.Engine {

	engine := gin.Default()

	engine.Use(CORSMiddleware(origin))

	api := engine.Group("/groups")
	api.Use(limits.RequestSizeLimiter(server.MaxBodyBytes))
	api.Use(func(ctx *gin.Context) {
		if len(ctx.Errors) > 0 {
			return
		}
	})

	apiAuth := api.Use(auth.MustAuthWithKey(server.PublicKey))

	apiAuth.GET("/", server.GetUserGroups)
	apiAuth.POST("/", server.CreateGroup)
	apiAuth.DELETE("/:groupID", server.DeleteGroup)

	apiAuth.POST("/:groupID/image", server.SetGroupProfilePicture)
	apiAuth.DELETE("/:groupID/image", server.DeleteGroupProfilePicture)

	apiAuth.DELETE("/:groupID/member/:memberID", server.DeleteUserFromGroup)
	apiAuth.PATCH("/:groupID/member/:memberID", server.GrantPriv)

	apiAuth.GET("/invites", server.GetUserInvites)
	apiAuth.POST("/invites", server.CreateInvite)
	apiAuth.PUT("/invites/:inviteID", server.RespondGroupInvite)

	return engine
}

func CORSMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
