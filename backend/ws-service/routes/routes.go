package routes

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/handlers"
	"github.com/gin-gonic/gin"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
)

func Setup(server *handlers.Server, origin string) *gin.Engine {
	engine := gin.Default()

	engine.Use(CORSMiddleware(origin))

	api := engine.Group("/ws")
	api.Use(auth.MustAuthWithKey(server.PublicKey))
	api.GET("/accessCode", server.GetAuthCode)

	engine.GET("/ws/websocket", server.ServeWebSocket)

	return engine
}

func CORSMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
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
