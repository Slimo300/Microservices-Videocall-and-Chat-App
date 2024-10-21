package models_test

import (
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"gotest.tools/assert"
)

func TestMemberRights(t *testing.T) {
	group := models.CreateGroup(uuid.New(), "new group")
	creatorUser := group.Members()[0]
	deletingUser := group.AddMember(uuid.New(), models.WithDeletingMembers)
	deletingUserTest := group.AddMember(uuid.New(), models.WithDeletingMembers)
	basicUser := group.AddMember(uuid.New())
	basicUserTest := group.AddMember(uuid.New())

	testCases := []struct {
		desc      string
		issuer    models.Member
		target    models.Member
		canDelete bool
		canAlter  bool
	}{
		{
			desc:      "creator_on_himself",
			issuer:    creatorUser,
			target:    creatorUser,
			canDelete: false,
			canAlter:  false,
		},
		{
			desc:      "creator_on_deleter",
			issuer:    creatorUser,
			target:    deletingUserTest,
			canDelete: true,
			canAlter:  true,
		},
		{
			desc:      "creator_on_basic",
			issuer:    creatorUser,
			target:    basicUserTest,
			canDelete: true,
			canAlter:  true,
		},
		{
			desc:      "deleter_on_creator",
			issuer:    deletingUser,
			target:    creatorUser,
			canDelete: false,
			canAlter:  false,
		},
		{
			desc:      "deleter_on_deleter",
			issuer:    deletingUser,
			target:    deletingUserTest,
			canDelete: false,
			canAlter:  false,
		},
		{
			desc:      "deleter_on_basic",
			issuer:    deletingUser,
			target:    basicUserTest,
			canDelete: true,
			canAlter:  false,
		},
		{
			desc:      "deleter_on_himself",
			issuer:    deletingUser,
			target:    deletingUser,
			canDelete: true,
			canAlter:  false,
		},
		{
			desc:      "basic_on_creator",
			issuer:    basicUser,
			target:    creatorUser,
			canDelete: false,
			canAlter:  false,
		},
		{
			desc:      "basic_on_deleter",
			issuer:    basicUser,
			target:    deletingUserTest,
			canDelete: false,
			canAlter:  false,
		},
		{
			desc:      "basic_on_basic",
			issuer:    basicUser,
			target:    basicUserTest,
			canDelete: false,
			canAlter:  false,
		},
		{
			desc:      "basic_on_himself",
			issuer:    basicUser,
			target:    basicUser,
			canDelete: true,
			canAlter:  false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.issuer.CanDelete(tC.target), tC.canDelete)
			assert.Equal(t, tC.issuer.CanAlter(tC.target), tC.canAlter)
		})
	}
}
