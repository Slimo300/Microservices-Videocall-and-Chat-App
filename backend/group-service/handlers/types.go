package handlers

import (
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
)

type group struct {
	ID         string   `json:"ID"`
	Name       string   `json:"name"`
	HasPicture bool     `json:"hasPicture"`
	Members    []member `json:"members"`
}

type member struct {
	ID               string `json:"ID"`
	GroupID          string `json:"groupID"`
	UserID           string `json:"userID"`
	User             user   `json:"user"`
	Adding           bool   `json:"adding"`
	DeletingMembers  bool   `json:"deletingMembers"`
	DeletingMessages bool   `json:"deletingMessages"`
	Muting           bool   `json:"muting"`
	Admin            bool   `json:"admin"`
	Creator          bool   `json:"creator"`
}

type user struct {
	ID         string `json:"ID"`
	Username   string `json:"username"`
	HasPicture bool   `json:"hasPicture"`
}

type invite struct {
	ID       string    `json:"ID"`
	IssuerID string    `json:"issuerID"`
	Issuer   user      `json:"user"`
	TargetID string    `json:"targetID"`
	Target   user      `json:"target"`
	GroupID  string    `json:"groupID"`
	Group    group     `json:"group"`
	Status   int       `json:"status"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func modelGroupToResponse(model models.Group) group {
	members := []member{}
	for _, mem := range model.Members() {
		members = append(members, modelMemberToResponse(mem))
	}

	return group{
		ID:         model.ID().String(),
		Name:       model.Name(),
		HasPicture: model.HasPicture(),
		Members:    members,
	}
}

func modelMemberToResponse(model models.Member) member {
	return member{
		ID:               model.ID().String(),
		UserID:           model.UserID().String(),
		GroupID:          model.GroupID().String(),
		User:             modelUserToResponse(model.User()),
		Adding:           model.Adding(),
		DeletingMembers:  model.DeletingMembers(),
		DeletingMessages: model.DeletingMessages(),
		Muting:           model.Muting(),
		Admin:            model.Admin(),
		Creator:          model.Creator(),
	}
}

func modelUserToResponse(model models.User) user {
	return user{
		ID:         model.ID().String(),
		Username:   model.Username(),
		HasPicture: model.HasPicture(),
	}
}

func modelInviteToResponse(model models.Invite) invite {
	return invite{
		ID:       model.ID().String(),
		GroupID:  model.GroupID().String(),
		Group:    modelGroupToResponse(model.Group()),
		TargetID: model.TargetID().String(),
		Target:   modelUserToResponse(model.Target()),
		IssuerID: model.IssuerID().String(),
		Issuer:   modelUserToResponse(model.Issuer()),
		Status:   int(model.Status()),
		Created:  model.Created(),
		Modified: model.Created(),
	}
}

func modelInvitesToResponse(modelInvites []models.Invite) []invite {
	is := []invite{}
	for _, i := range modelInvites {
		is = append(is, modelInviteToResponse(i))
	}
	return is
}

func modelGroupsToResponse(modelGroups []models.Group) []group {
	gs := []group{}
	for _, g := range modelGroups {
		gs = append(gs, modelGroupToResponse(g))
	}
	return gs
}
