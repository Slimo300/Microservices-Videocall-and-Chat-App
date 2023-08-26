package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	w "github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/webrtc"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) ServeWebSocket(c *gin.Context) {

	reqMemberID := c.Param("memberID")
	if reqMemberID == "" {
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

	if memberID != reqMemberID {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid access code"})
		return
	}

	member, err := s.DB.GetMember(memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
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
