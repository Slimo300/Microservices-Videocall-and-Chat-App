package events

import "github.com/google/uuid"

type MemberUpdatedEvent struct {
	ID               uuid.UUID `json:"id"`
	DeletingMessages *bool     `json:"deletingMessages" binding:"required"`
	DeletingMembers  *bool     `json:"deletingMembers" binding:"required"`
	Adding           *bool     `json:"adding" binding:"required"`
	Setting          *bool     `json:"setting" binding:"required"`
}

func (MemberUpdatedEvent) EventName() string {
	return "member.updated"
}
