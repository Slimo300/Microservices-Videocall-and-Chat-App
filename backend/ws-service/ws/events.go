package ws

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

// Deletes group from every user that is subscribed to it and sends information via websocket to user
func (h *WSHub) GroupDeleted(event events.GroupDeletedEvent) {
	for client := range h.clients {
		for i, group := range client.groups {
			if group == event.ID {
				client.groups = append(client.groups[:i], client.groups[:i+1]...)
				client.send <- &Action{ActionType: "DELETE_GROUP", Payload: event}
			}
		}
	}
}

// Adds subscription to member groups and sends info to other members in group
func (h *WSHub) MemberAdded(event events.MemberCreatedEvent) {
	for client := range h.clients {
		if client.id == event.UserID {
			client.groups = append(client.groups, event.GroupID)
			continue
		}
		if !event.Creator {
			for _, group := range client.groups {
				if event.GroupID == group {
					client.send <- &Action{ActionType: "ADD_MEMBER", Payload: event}
				}
			}
		}
	}
}

// Deletes member subscription and sends info about it to other members in group
func (h *WSHub) MemberDeleted(event events.MemberDeletedEvent) {
	for client := range h.clients {
		for i, group := range client.groups {
			// if user is a member of group
			if group == event.GroupID {
				// if user is the one to be deleted
				if client.id == event.UserID {
					client.groups = append(client.groups[:i], client.groups[:i+1]...)
				} else {
					client.send <- &Action{ActionType: "DELETE_MEMBER", Payload: event}
				}
			}
		}
	}
}

func (h *WSHub) MemberUpdated(event events.MemberUpdatedEvent) {
	for client := range h.clients {
		for _, group := range client.groups {
			if group == event.GroupID {
				if client.id == event.UserID {
					client.send <- &Action{ActionType: "UPDATE_MEMBER", Payload: event}
				}
			}
		}
	}
}

// Sends invite to specified user
func (h *WSHub) InviteSent(event events.InviteSentEvent) {
	for client := range h.clients {
		if client.id == event.TargetID {
			client.send <- &Action{ActionType: "SEND_INVITE", Payload: event}
		}
	}
}

func (h *WSHub) MessageDeleted(event events.MessageDeletedEvent) {
	for client := range h.clients {
		for _, group := range client.groups {
			if group == event.GroupID {
				client.send <- &Action{ActionType: "DELETE_MESSAGE", Payload: event}
			}
		}
	}
}
