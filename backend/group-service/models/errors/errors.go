package models

import "fmt"

type errInviteNotFound struct {
	inviteID string
}

func (e errInviteNotFound) Error() string {
	return fmt.Sprintf("invite with id %s not found", e.inviteID)
}

func NewInviteNotFoundError(inviteID string) error {
	return errInviteNotFound{inviteID: inviteID}
}

type errUserNotFound struct {
	userID string
}

func (e errUserNotFound) Error() string {
	return fmt.Sprintf("user with id %s not found", e.userID)
}

func NewUserNotFoundError(userID string) error {
	return errUserNotFound{userID: userID}
}

type errMemberNotFound struct {
	memberID string
}

func (e errMemberNotFound) Error() string {
	return fmt.Sprintf("member with id %s not found", e.memberID)
}

func NewMemberNotFoundError(memberID string) error {
	return errMemberNotFound{memberID: memberID}
}

type errUserNotInGroup struct {
	userID, groupID string
}

func (e errUserNotInGroup) Error() string {
	return fmt.Sprintf("user with id %s is not a member of group %s", e.userID, e.groupID)
}

func NewUserNotInGroupError(userID, groupID string) error {
	return errUserNotInGroup{userID: userID, groupID: groupID}
}

type errUserAlreadyInGroup struct {
	userID, groupID string
}

func (e errUserAlreadyInGroup) Error() string {
	return fmt.Sprintf("user with id %s is already a member of group %s", e.userID, e.groupID)
}

func NewUserAlreadyInGroupError(userID, groupID string) error {
	return errUserAlreadyInGroup{userID: userID, groupID: groupID}
}

type errUserAlreadyInvited struct {
	userID, groupID string
}

func (e errUserAlreadyInvited) Error() string {
	return fmt.Sprintf("user with id %s already invited to group %s", e.userID, e.groupID)
}

func NewUserAlreadyInvitedError(userID, groupID string) error {
	return errUserAlreadyInvited{userID: userID, groupID: groupID}
}

type errUserIsNotInvitesTarget struct {
	userID, inviteID string
}

func (e errUserIsNotInvitesTarget) Error() string {
	return fmt.Sprintf("user with id %s is not a target of invite %s", e.userID, e.inviteID)
}

func NewUserIsNotInvitesTargetError(userID, inviteID string) error {
	return errUserIsNotInvitesTarget{userID: userID, inviteID: inviteID}
}

type errInviteAlreadyAnswered struct {
	inviteID string
}

func (e errInviteAlreadyAnswered) Error() string {
	return fmt.Sprintf("invite with id %s already answered", e.inviteID)
}

func NewInviteAlreadyAnsweredError(inviteID string) error {
	return errInviteAlreadyAnswered{inviteID: inviteID}
}

type Action struct{ a string }

func DeleteMemberAction() Action { return Action{a: "DELETE_MEMBER"} }
func AddMemberAction() Action    { return Action{a: "ADD_MEMBER"} }
func UpdateMemberAction() Action { return Action{a: "UPDATE_MEMBER"} }
func UpdateGroupAction() Action  { return Action{a: "UPDATE_GROUP"} }

type errMemberUnauthorized struct {
	groupID string
	action  Action
}

func (e errMemberUnauthorized) Error() string {
	return fmt.Sprintf("member has no rights to %s on group %s", e.action, e.groupID)
}

func NewMemberUnauthorizedError(groupID string, action Action) error {
	return errMemberUnauthorized{groupID: groupID, action: action}
}

type errUserCantSeeInvite struct {
	userID   string
	inviteID string
}

func (e errUserCantSeeInvite) Error() string {
	return fmt.Sprintf("user %s can't see invite %s", e.userID, e.inviteID)
}

func NewUserCantSeeInviteError(userID, inviteID string) error {
	return errUserCantSeeInvite{userID: userID, inviteID: inviteID}
}
