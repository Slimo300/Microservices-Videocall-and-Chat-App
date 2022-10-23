package routes

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/handlers"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, server *handlers.Server) {

	engine.Use(CORSMiddleware())

	api := engine.Group("/api")
	api.Use(server.CheckDatabase())
	api.Use(limits.RequestSizeLimiter(server.MaxBodyBytes))
	apiAuth := api.Use(server.MustAuth())

	apiAuth.GET("/group", server.GetUserGroups)
	apiAuth.POST("/group", server.CreateGroup)
	apiAuth.DELETE("/group/:groupID", server.DeleteGroup)
	apiAuth.POST("/group/:groupID/image", server.SetGroupProfilePicture)
	apiAuth.DELETE("/group/:groupID/image", server.DeleteGroupProfilePicture)

	apiAuth.DELETE("/member/:memberID", server.DeleteUserFromGroup)
	apiAuth.PUT("/member/:memberID", server.GrantPriv)

	apiAuth.GET("/invites", server.GetUserInvites)
	apiAuth.POST("/invites", server.SendGroupInvite)
	apiAuth.PUT("/invites/:inviteID", server.RespondGroupInvite)

	ws := engine.Group("/ws")
	ws.Use(server.AuthWS())
	ws.GET("", server.ServeWebSocket)
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
