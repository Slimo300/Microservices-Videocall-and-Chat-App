package routes

import (
	"github.com/Slimo300/chat-wsservice/internal/handlers"
	"github.com/gin-gonic/gin"

	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
)

func Setup(server *handlers.Server, origin string) *gin.Engine {
	engine := gin.Default()

	engine.Use(CORSMiddleware(origin))

	api := engine.Group("/ws")
	api.Use(tokens.MustAuth(server.TokenClient))
	api.GET("/accessCode", server.GetAuthCode)

	engine.GET("/ws", server.ServeWebSocket)

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
