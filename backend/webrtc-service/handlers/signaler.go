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

	accessCode := c.Query("accessCode")
	if accessCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no access code provided"})
		return
	}

	reqGroupID := c.Param("groupID")
	if reqGroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no group id provided"})
		return
	}

	userID, groupID, err := s.DB.CheckAccessCode(accessCode)
	if err != nil || groupID != reqGroupID {
		c.JSON(http.StatusBadRequest, gin.H{"err": "connection not authorized"})
		return
	}

	username, err := s.DB.GetMember(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	streamID := c.Query("streamID")
	if streamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "streamID not specified"})
		return
	}

	enabled := true
	notEnabled := false

	userData := w.UserConnData{
		Username:     username,
		StreamID:     streamID,
		AudioEnabled: &enabled,
	}

	videoEnabled := c.Query("videoEnabled")
	if videoEnabled == "true" {
		userData.VideoEnabled = &enabled
	} else {
		userData.VideoEnabled = &notEnabled
	}

	room := s.Relay.GetRoom(groupID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	room.ConnectRoom(conn, userData)
}
