package routes

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(server *handlers.Server, origin string) *gin.Engine {
	engine := gin.Default()

	engine.Use(CORSMiddleware(origin))

	api := engine.Group("/messages")
	apiAuth := api.Use(auth.MustAuthWithKey(server.PublicKey))

	apiAuth.GET("/group/:groupID/messages", server.GetGroupMessages)
	apiAuth.DELETE("/group/:groupID/messages/:messageID", server.DeleteMessageForEveryone)
	apiAuth.PATCH("/group/:groupID/messages/:messageID", server.DeleteMessageForYourself)
	apiAuth.POST("/group/:groupID", server.GetPresignedPutRequest)

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
