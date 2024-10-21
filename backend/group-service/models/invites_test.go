package models_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAnswerInvite(t *testing.T) {
	groupID := uuid.New()
	issuerID := uuid.New()
	targetID := uuid.New()

	inviteToDecline := models.CreateInvite(issuerID, targetID, groupID)
	inviteToAccept := models.CreateInvite(issuerID, targetID, groupID)

	testCases := []struct {
		desc         string
		invite       *models.Invite // !IMPORTANT here we pass a pointer to invite so the same invite might be used many times, this should be changed if parallelism would be introduced
		userID       uuid.UUID
		answer       bool
		checkResults func(t *testing.T, invite models.Invite, err error)
	}{
		{
			desc:   "user_not_a_target",
			invite: &inviteToDecline,
			userID: issuerID,
			answer: false,
			checkResults: func(t *testing.T, invite models.Invite, err error) {
				assert.Error(t, err, fmt.Sprintf("user with id %s is not a target of invite %s", issuerID, invite.ID().String()))
			},
		},
		{
			desc:   "invite_declined",
			invite: &inviteToDecline,
			userID: targetID,
			answer: false,
			checkResults: func(t *testing.T, invite models.Invite, err error) {
				assert.NoError(t, err)
				assert.Equal(t, invite.Status(), models.INVITE_DECLINED)
			},
		},
		{
			desc:   "invite_accepted",
			invite: &inviteToAccept,
			userID: targetID,
			answer: true,
			checkResults: func(t *testing.T, invite models.Invite, err error) {
				assert.NoError(t, err)
				assert.Equal(t, invite.Status(), models.INVITE_ACCEPTED)
			},
		},
		{
			desc:   "invite_answered",
			invite: &inviteToAccept,
			userID: targetID,
			answer: false,
			checkResults: func(t *testing.T, invite models.Invite, err error) {
				assert.Error(t, err, fmt.Sprintf("invite with id %s already answered", invite.ID().String()))
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			log.Println(tC.invite.Status())
			err := tC.invite.AnswerInvite(tC.userID, tC.answer)
			tC.checkResults(t, *tC.invite, err)
		})
	}
}
