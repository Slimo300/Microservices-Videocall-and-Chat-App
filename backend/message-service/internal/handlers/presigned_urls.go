package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type presignedEnvelope struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

func (s *Server) GetPresignedPutRequest(c *gin.Context) {
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

	files := c.Query("files")
	if files == "" {
		files = "1"
	}
	filesInt, err := strconv.Atoi(files)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "files query value is not a valid integer"})
		return
	}

	if _, err := s.DB.GetGroupMembership(userID, groupID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "user cannot send messages to this group"})
		return
	}

	var requestsData []presignedEnvelope
	for i := 0; i < filesInt; i++ {
		key := groupID.String() + "/" + uuid.NewString()
		url, err := s.Storage.GetPresignedPutRequest(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		requestsData = append(requestsData, presignedEnvelope{Url: url, Key: key})
	}

	c.JSON(http.StatusOK, gin.H{"requests": requestsData})

}
