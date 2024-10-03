package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) SearchUsers(c *gin.Context) {
	username := c.Param("name")
	if !IsValidUsername(username) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "not a valid username"})
		return
	}
	num, err := strconv.Atoi(c.Query("num"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "num is not a valid numer"})
		return
	}
	result, err := s.DB.GetUsers(username, num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// TODO REGEXP
func IsValidUsername(name string) bool {
	return true
}
