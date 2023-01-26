package handlers

import (
	"net/http"
	"strconv"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetGroupMessages(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	num, err := strconv.Atoi(c.Query("num"))
	if err != nil || num <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "number of messages is not a valid number"})
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "offset is not a valid number"})
		return
	}

	messages, err := s.DB.GetGroupMessages(userID, groupID, offset, num)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}
	if len(messages) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (s *Server) DeleteMessageForEveryone(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid user ID"})
		return
	}
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}
	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid message ID"})
		return
	}

	msg, err := s.DB.DeleteMessageForEveryone(userID, messageID, groupID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	s.Emitter.Emit(events.MessageDeletedEvent{
		ID:      msg.ID,
		GroupID: msg.GroupID,
	})

	c.JSON(http.StatusOK, msg)
}

func (s *Server) DeleteMessageForYourself(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid user ID"})
		return
	}
	groupID, err := uuid.Parse((c.Param("groupID")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}

	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid message ID"})
		return
	}

	msg, err := s.DB.DeleteMessageForYourself(userID, messageID, groupID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, msg)

}

func (s *Server) GetPresignedPutRequest(c *gin.Context) {
	// userID, err := uuid.Parse(c.GetString("userID"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"err": "invalid user ID"})
	// 	return
	// }
	// groupID, err := uuid.Parse((c.Param("groupID")))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
	// 	return
	// }
	// files := c.Query("files")
	// if files == "" {
	// 	files = "1"
	// }
	// filesInt, err := strconv.Atoi(files)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"err": "files query value is not a valid integer"})
	// 	return
	// }

	// check Membership
	// generate presigned urls
}
