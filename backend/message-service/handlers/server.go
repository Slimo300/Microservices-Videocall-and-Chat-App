package handlers

import (
	"crypto/rsa"
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/app"
	"github.com/gin-gonic/gin"
)

type Server struct {
	App       app.App
	PublicKey *rsa.PublicKey
}

func NewServer(app app.App, pubKey *rsa.PublicKey, origin string) http.Handler {
	s := &Server{
		App:       app,
		PublicKey: pubKey,
	}
	return s.setup(origin)
}

func (s *Server) setup(origin string) http.Handler {
	engine := gin.Default()

	engine.Use(CORSMiddleware(origin))

	api := engine.Group("/messages")
	apiAuth := api.Use(auth.MustAuthWithKey(s.PublicKey))

	apiAuth.GET("/:groupID", s.GetGroupMessages)
	apiAuth.DELETE("/:messageID", s.DeleteMessageForEveryone)
	apiAuth.PATCH("/:messageID/hide", s.DeleteMessageForYourself)
	apiAuth.POST("/:groupID/presign/put", s.GetPresignedPutRequests)
	apiAuth.POST("/:groupID/presign/get", s.GetPresignedGetRequests)

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
