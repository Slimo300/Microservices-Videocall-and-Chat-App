package handlers

import (
	"crypto/rsa"
	"net/http"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

type Server struct {
	App            app.App
	PublicKey      *rsa.PublicKey
	MaxBodyBytes   int64
	ServiceAddress string
}

func NewServer(a app.App, pubKey *rsa.PublicKey, origin, serviceAddress string) http.Handler {
	s := Server{
		App:            a,
		MaxBodyBytes:   4194304,
		PublicKey:      pubKey,
		ServiceAddress: serviceAddress,
	}

	return s.setup(origin)
}

func (s *Server) setup(origin string) http.Handler {
	engine := gin.Default()

	engine.Use(CORSMiddleware(origin))

	api := engine.Group("/groups")
	api.Use(limits.RequestSizeLimiter(s.MaxBodyBytes))
	api.Use(func(ctx *gin.Context) {
		if len(ctx.Errors) > 0 {
			return
		}
	})

	apiAuth := api.Use(auth.MustAuthWithKey(s.PublicKey))

	apiAuth.GET("/", s.GetUserGroups)
	apiAuth.POST("/", s.CreateGroup)
	apiAuth.GET("/:groupID", s.GetGroupByID)
	apiAuth.DELETE("/:groupID", s.DeleteGroup)

	apiAuth.POST("/:groupID/image", s.SetGroupProfilePicture)
	apiAuth.DELETE("/:groupID/image", s.DeleteGroupProfilePicture)

	apiAuth.DELETE("/members/:memberID", s.DeleteUserFromGroup)
	apiAuth.PATCH("/members/:memberID", s.GrantRights)

	apiAuth.GET("/invites", s.GetUserInvites)
	apiAuth.POST("/invites", s.CreateInvite)
	apiAuth.GET("/invites/:inviteID", s.GetInviteByID)
	apiAuth.PUT("/invites/:inviteID", s.RespondGroupInvite)

	return engine
}

func CORSMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Location")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
