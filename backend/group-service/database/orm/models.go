package orm

import (
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

type Group struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Name       string
	HasPicture bool
	Members    []Member `gorm:"foreignKey:GroupID"`
}

func marshalGroup(group models.Group) Group {
	members := []Member{}
	for _, mem := range group.Members() {
		members = append(members, marshalMember(mem))
	}
	return Group{
		ID:         group.ID(),
		Name:       group.Name(),
		HasPicture: group.HasPicture(),
		Members:    members,
	}
}

func unmarshalGroup(group Group) models.Group {
	var members []models.Member
	for _, m := range group.Members {
		members = append(members, unmarshalMember(m))
	}
	return models.UnmarshalGroupFromDatabase(group.ID, group.Name, group.HasPicture, members)
}

func unmarshalInvites(invites []Invite) []models.Invite {
	var is []models.Invite
	for _, i := range invites {
		is = append(is, unmarshalInvite(i))
	}
	return is
}

func unmarshalGroups(groups []Group) []models.Group {
	var gs []models.Group
	for _, g := range groups {
		gs = append(gs, unmarshalGroup(g))
	}
	return gs
}

func (Group) TableName() string { return "groups" }

type Member struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	GroupID          uuid.UUID `gorm:"uniqueIndex:idx_first;size:191"`
	Group            Group     `gorm:"constraint:OnDelete:CASCADE;"`
	UserID           uuid.UUID `gorm:"uniqueIndex:idx_first;size:191"`
	User             User      `gorm:"constraint:OnDelete:CASCADE;"`
	Adding           *bool
	DeletingMembers  *bool
	DeletingMessages *bool
	Muting           *bool
	Admin            *bool
	Creator          *bool
}

func (Member) TableName() string { return "members" }

func marshalMember(member models.Member) Member {
	adding := member.Adding()
	deletingMembers := member.DeletingMembers()
	deletingMessages := member.DeletingMessages()
	muting := member.Muting()
	admin := member.Admin()
	creator := member.Creator()

	return Member{
		ID:               member.ID(),
		GroupID:          member.GroupID(),
		UserID:           member.UserID(),
		Adding:           &adding,
		DeletingMembers:  &deletingMembers,
		DeletingMessages: &deletingMessages,
		Muting:           &muting,
		Admin:            &admin,
		Creator:          &creator,
	}
}
func unmarshalMember(m Member) models.Member {
	return models.UnmarshalMemberFromDatabase(
		m.ID,
		m.UserID,
		m.GroupID,
		models.UnmarshalUserFromDatabase(m.User.ID, m.User.UserName, m.User.HasPicture),
		*m.Adding,
		*m.DeletingMembers,
		*m.DeletingMessages,
		*m.Muting,
		*m.Admin,
		*m.Creator,
	)
}

type User struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	UserName   string    `gorm:"unique"`
	HasPicture bool
}

func (User) TableName() string { return "users" }

func marshalUser(user models.User) User {
	return User{
		ID:         user.ID(),
		HasPicture: user.HasPicture(),
		UserName:   user.Username(),
	}
}

type Invite struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	IssId    uuid.UUID `gorm:"size:191"`
	Iss      User      `gorm:"constraint:OnDelete:CASCADE"`
	TargetID uuid.UUID `gorm:"size:191"`
	Target   User      `gorm:"constraint:OnDelete:CASCADE"`
	GroupID  uuid.UUID `gorm:"size:191"`
	Group    Group     `gorm:"constraint:OnDelete:CASCADE"`
	Status   int
	Created  time.Time
	Modified time.Time
}

func (Invite) TableName() string { return "invites" }

func marshalInvite(invite models.Invite) Invite {
	return Invite{
		ID:       invite.ID(),
		GroupID:  invite.GroupID(),
		IssId:    invite.IssuerID(),
		TargetID: invite.TargetID(),
		Status:   int(invite.Status()),
		Created:  invite.Created(),
		Modified: invite.Modified(),
	}
}

func unmarshalInvite(invite Invite) models.Invite {
	group := unmarshalGroup(invite.Group)
	issuer := models.UnmarshalUserFromDatabase(invite.Iss.ID, invite.Iss.UserName, invite.Iss.HasPicture)
	target := models.UnmarshalUserFromDatabase(invite.Target.ID, invite.Target.UserName, invite.Target.HasPicture)
	return models.UnmarshalInviteFromDatabase(
		invite.ID,
		invite.GroupID,
		invite.IssId,
		invite.TargetID,
		group,
		issuer,
		target,
		models.InviteStatus(invite.Status),
		invite.Created,
		invite.Modified,
	)
}
