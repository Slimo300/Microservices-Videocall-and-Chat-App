package handlers

import (
	"net/http"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/gin-gonic/gin"
)

func (server *Server) Setup(origin string) http.Handler {
	engine := gin.Default()

	engine.Use(CORSMiddleware(origin))

	api := engine.Group("/video-call")
	api.Use(auth.MustAuthWithKey(server.PublicKey))
	api.GET("/:groupID/accessCode", server.GetAuthCode)

	engine.GET("/video-call/:groupID/ws", server.ServeWebSocket)

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
