package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) SearchUser(c *gin.Context) {
	username := c.Param("name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid name"})
		return
	}

	user, err := s.DB.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": fmt.Sprintf("user with name %s not found", username)})
		return
	}

	c.JSON(http.StatusOK, user)
}
