package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	room := s.Relay.GetRoom(groupID)
	// check if room for group exists
	// room, ok := s.Relay[groupID]
	// if !ok {
	// 	room = &w.Room{}
	// 	room.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	// 	s.Relay[groupID] = room
	// }

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal server error"})
		return
	}

	room.ConnectRoom(conn, username)
}
