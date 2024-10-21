package models_test

import (
	"log"
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangePictureIfStateIncorrect(t *testing.T) {
	group := models.CreateGroup(uuid.New(), "new group")

	require.False(t, group.HasPicture())

	assert.True(t, group.ChangePictureStateIfIncorrect(true))
	assert.True(t, group.HasPicture())

	assert.False(t, group.ChangePictureStateIfIncorrect(true))
	assert.True(t, group.HasPicture())

	assert.True(t, group.ChangePictureStateIfIncorrect(false))
	assert.False(t, group.HasPicture())

	assert.False(t, group.ChangePictureStateIfIncorrect(false))
	assert.False(t, group.HasPicture())
}

func TestGetMemberByID(t *testing.T) {
	group := models.CreateGroup(uuid.New(), "new group")

	userID := uuid.New()

	returnedMember, ok := group.GetMemberByUserID(userID)
	assert.False(t, ok)
	assert.Empty(t, returnedMember)

	group.AddMember(userID)

	returnedMember, ok = group.GetMemberByUserID(userID)
	assert.True(t, ok)
	assert.NotEmpty(t, returnedMember)
	assert.Len(t, group.Members(), 2)
	log.Println(group.Members())
}
