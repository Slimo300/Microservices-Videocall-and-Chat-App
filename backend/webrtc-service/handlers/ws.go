package handlers

import (
	"log"
	"net/http"

	w "github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/webrtc"
	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
)

func (s *Server) ServeWebSocket(c *gin.Context) {

	accessCode := c.Query("accessCode")
	if accessCode == "" {
		log.Println("Access Code Invalid")
		c.JSON(http.StatusBadRequest, gin.H{"err": "no access code provided"})
		return
	}

	reqGroupID := c.Param("groupID")
	if reqGroupID == "" {
		log.Println("reqGroupID invalid")
		c.JSON(http.StatusBadRequest, gin.H{"err": "no group id provided"})
		return
	}

	groupID, err := s.DB.CheckAccessCode(accessCode)
	if err != nil || groupID != reqGroupID {
		log.Printf("ReqGroupID: %s, groupID: %s, Err: %v", reqGroupID, groupID, err)
		c.JSON(http.StatusBadRequest, gin.H{"err": "connection not authorized"})
		return
	}

	// check if room for group exists
	room, ok := s.Relay[groupID]
	if !ok {
		room = &w.Room{}
		room.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
		s.Relay[groupID] = room
	}

	room.ConnectRoom(c.Writer, c.Request)
}
