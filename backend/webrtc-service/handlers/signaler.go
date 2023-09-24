package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	w "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/webrtc"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) ServeWebSocket(c *gin.Context) {

	reqGroupID := c.Param("groupID")
	if reqGroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no group id provided"})
		return
	}

	accessCode := c.Query("accessCode")
	if accessCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no access code provided"})
		return
	}

	streamID := c.Query("streamID")
	if streamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "streamID not specified"})
		return
	}

	memberID, err := s.DB.CheckAccessCode(accessCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "connection not authorized"})
		return
	}

	member, err := s.DB.GetMemberByID(memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	if member.GroupID != reqGroupID {
		c.JSON(http.StatusForbidden, gin.H{"err": "Supplied access code is for a wrong group"})
		return
	}

	userData := w.UserConnData{
		Username:     member.Username,
		StreamID:     streamID,
		AudioEnabled: applyMediaQuery(c.Query("audio")),
		VideoEnabled: applyMediaQuery(c.Query("video")),
	}

	room := s.Relay.GetRoom(member.GroupID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	room.ConnectRoom(conn, userData)
}

func applyMediaQuery(query string) *bool {

	enabled := true
	notEnabled := false

	if query == "true" {
		return &enabled
	}
	return &notEnabled
}
