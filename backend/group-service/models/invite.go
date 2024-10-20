package models

import (
	"time"

	"github.com/google/uuid"

	merrors "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models/errors"
)

type InviteStatus int

const (
	INVITE_AWAITING InviteStatus = iota
	INVITE_ACCEPTED
	INVITE_DECLINED
)

type Invite struct {
	id       uuid.UUID
	issuerID uuid.UUID
	issuer   User
	targetID uuid.UUID
	target   User
	groupID  uuid.UUID
	group    Group
	status   InviteStatus
	created  time.Time
	modified time.Time
}

func (i Invite) ID() uuid.UUID        { return i.id }
func (i Invite) GroupID() uuid.UUID   { return i.groupID }
func (i Invite) IssuerID() uuid.UUID  { return i.issuerID }
func (i Invite) TargetID() uuid.UUID  { return i.targetID }
func (i Invite) Group() Group         { return i.group }
func (i Invite) Issuer() User         { return i.issuer }
func (i Invite) Target() User         { return i.target }
func (i Invite) Status() InviteStatus { return i.status }
func (i Invite) Modified() time.Time  { return i.modified }
func (i Invite) Created() time.Time   { return i.created }

func (i Invite) CanUserSee(userID uuid.UUID) bool {
	return userID == i.issuerID || userID == i.targetID
}

func (i Invite) CanAnswer(userID uuid.UUID) error {
	if i.targetID != userID {
		return merrors.NewUserIsNotInvitesTargetError(userID.String(), i.id.String())
	}
	if i.status != INVITE_AWAITING {
		return merrors.NewInviteAlreadyAnsweredError(i.id.String())
	}
	return nil
}

func (i *Invite) AnswerInvite(userID uuid.UUID, answer bool) error {
	if err := i.CanAnswer(userID); err != nil {
		return err
	}
	i.modified = time.Now()
	if answer {
		i.status = INVITE_ACCEPTED
	} else {
		i.status = INVITE_DECLINED
	}
	return nil
}

func CreateInvite(issuerID, targetID, groupID uuid.UUID) Invite {
	return Invite{
		id:       uuid.New(),
		issuerID: issuerID,
		targetID: targetID,
		groupID:  groupID,
		status:   INVITE_AWAITING,
		created:  time.Now(),
		modified: time.Now(),
	}
}

func UnmarshalInviteFromDatabase(inviteID, groupID, issuerID, targetID uuid.UUID, group Group, issuer, target User, status InviteStatus, created, modified time.Time) Invite {
	return Invite{
		id:       inviteID,
		groupID:  groupID,
		issuerID: issuerID,
		targetID: targetID,
		group:    group,
		issuer:   issuer,
		target:   target,
		status:   status,
		created:  created,
		modified: modified,
	}
}
