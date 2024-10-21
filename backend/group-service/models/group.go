package models

import (
	"github.com/google/uuid"
)

type Group struct {
	id         uuid.UUID
	name       string
	hasPicture bool
	members    []Member
}

func (g Group) ID() uuid.UUID     { return g.id }
func (g Group) Name() string      { return g.name }
func (g Group) HasPicture() bool  { return g.hasPicture }
func (g Group) Members() []Member { return g.members }

func (g *Group) ChangePictureStateIfIncorrect(state bool) bool {
	if g.hasPicture == state {
		return false
	}
	g.hasPicture = state
	return true
}

func (g *Group) AddMember(userID uuid.UUID, opts ...MemberOpt) Member {
	m := NewMember(userID, g.id)
	for _, o := range opts {
		o(&m)
	}
	g.members = append(g.members, m)
	return m
}

type MemberOpt func(m *Member)

func WithAdding(m *Member)          { m.adding = true }
func WithDeletingMembers(m *Member) { m.deletingMembers = true }

func (g Group) GetMemberByUserID(userID uuid.UUID) (Member, bool) {
	for _, m := range g.members {
		if m.userID == userID {
			return m, true
		}
	}
	return Member{}, false
}

func CreateGroup(userID uuid.UUID, name string) Group {
	groupID := uuid.New()
	creator := newCreatorMember(userID, groupID)
	return Group{
		id:         groupID,
		name:       name,
		hasPicture: false,
		members:    []Member{creator},
	}
}

// Only for unmarshaling from database
func UnmarshalGroupFromDatabase(groupID uuid.UUID, name string, hasPicture bool, members []Member) Group {
	return Group{
		id:         groupID,
		name:       name,
		hasPicture: hasPicture,
		members:    members,
	}
}
