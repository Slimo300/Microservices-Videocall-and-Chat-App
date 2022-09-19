package ws

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/communication"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/google/uuid"
)

// Group created updates user groups list to which he listens
func (h *Hub) GroupCreated(userID, groupID uuid.UUID) {
	for client := range h.clients {
		if client.id == userID {
			client.groups = append(client.groups, groupID)
		}
	}
}

// Deletes group from every user that is subscribed to it and sends information via websocket to user
func (h *Hub) GroupDeleted(groupID uuid.UUID) {
	for client := range h.clients {
		for i, group := range client.groups {
			if group == groupID {
				client.groups = append(client.groups[:i], client.groups[:i+1]...)
				client.send <- &communication.Action{Action: "DELETE_GROUP", Group: groupID}
			}
		}
	}
}

// Adds subscription to member groups and sends info to other members in group
func (h *Hub) MemberAdded(member models.Member) {
	for client := range h.clients {
		if client.id == member.UserID {
			client.groups = append(client.groups, member.GroupID)
			continue
		}
		for _, group := range client.groups {
			if member.GroupID == group {
				client.send <- &communication.Action{Action: "ADD_MEMBER", Member: member}
			}
		}
	}
}

// Deletes member subscription and sends info about it to other members in group
func (h *Hub) MemberDeleted(member models.Member) {
	for client := range h.clients {
		for i, group := range client.groups {
			// if user is a member of group
			if group == member.GroupID {
				// if user is the one to be deleted
				if client.id == member.UserID {
					client.groups = append(client.groups[:i], client.groups[:i+1]...)
				} else {
					client.send <- &communication.Action{Action: "DELETE_MEMBER", Member: member}
				}
			}
		}
	}
}

// Sends invite to specified user
func (h *Hub) SendGroupInvite(invite models.Invite) {
	for client := range h.clients {
		if client.id == invite.TargetID {
			client.send <- &communication.Action{Action: "SEND_INVITE", Invite: invite}
		}
	}
}
