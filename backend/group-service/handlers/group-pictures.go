package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) SetGroupProfilePicture(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid ID"})
		return
	}

	imageFileHeader, err := c.FormFile("avatarFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	mimeType := imageFileHeader.Header.Get("Content-Type")
	if !isAllowedImageType(mimeType) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "image extention not allowed"})
		return
	}

	groupID := c.Param("groupID")
	groupUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid group ID"})
		return
	}

	member, err := s.DB.GetUserGroupMember(userUID, groupUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if !member.Setting {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to set"})
		return
	}

	pictureURL, err := s.DB.GetGroupProfilePicture(groupUID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
		return
	}
	wasEmpty := false
	if pictureURL == "" {
		pictureURL = uuid.New().String()
		wasEmpty = true
	}

	file, err := imageFileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad image"})
		return
	}

	if err = s.Storage.UpdateProfilePicture(file, pictureURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if wasEmpty {
		if err = s.DB.SetGroupProfilePicture(groupUID, pictureURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"newUrl": pictureURL})
}

func (s *Server) DeleteGroupProfilePicture(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	groupID := c.Param("groupID")
	groupUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}
	member, err := s.DB.GetUserGroupMember(userUID, groupUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if member.Deleted {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to set"})
		return
	}
	if !member.Setting {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to set"})
		return
	}
	url, err := s.DB.GetGroupProfilePicture(groupUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "group has no image to delete"})
		return
	}
	if err = s.DB.SetGroupProfilePicture(groupUID, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if err = s.Storage.DeleteProfilePicture(url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
