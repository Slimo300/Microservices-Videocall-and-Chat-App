package ws

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

func (h *WSHub) messageSent(event events.MessageSentEvent) {

	if event.ServiceID == h.serviceID {
		return
	}

	var files []MessageFile
	for _, file := range event.Files {
		files = append(files, MessageFile{Key: file.Key, Ext: file.Extension})
	}

	for client := range h.clients {
		if _, ok := client.groups[event.GroupID]; ok {
			client.send <- &Message{
				ID:      event.ID,
				Group:   event.GroupID,
				User:    event.UserID,
				Nick:    event.Nick,
				Message: event.Text,
				When:    event.Posted,
				Files:   files,
			}
		}
	}
}

// Deletes group from every user that is subscribed to it and sends information via websocket to user
func (h *WSHub) groupDeleted(event events.GroupDeletedEvent) {
	for client := range h.clients {
		if _, ok := client.groups[event.ID]; ok {
			delete(client.groups, event.ID)
			client.send <- &Action{ActionType: "DELETE_GROUP", Payload: event.ID}
		}
	}
}

// Adds subscription to member groups and sends info to other members in group
func (h *WSHub) memberAdded(event events.MemberCreatedEvent) {
	for client := range h.clients {
		if client.id == event.UserID {
			client.groups[event.GroupID] = struct{}{}
			continue
		}
		if _, ok := client.groups[event.GroupID]; ok {
			client.send <- &Action{ActionType: "ADD_MEMBER", Payload: event}
		}
	}
}

// Deletes member subscription and sends info about it to other members in group
func (h *WSHub) memberDeleted(event events.MemberDeletedEvent) {
	for client := range h.clients {
		if _, ok := client.groups[event.GroupID]; ok {
			if client.id == event.UserID {
				delete(client.groups, event.GroupID)
				client.send <- &Action{ActionType: "DELETE_GROUP", Payload: event.GroupID}
			} else {
				client.send <- &Action{ActionType: "DELETE_MEMBER", Payload: event}
			}
		}
	}
}

func (h *WSHub) memberUpdated(event events.MemberUpdatedEvent) {
	for client := range h.clients {
		if _, ok := client.groups[event.GroupID]; ok {
			client.send <- &Action{ActionType: "UPDATE_MEMBER", Payload: event}
		}
	}
}

// Sends invite to specified user
func (h *WSHub) inviteSent(event events.InviteSentEvent) {
	for client := range h.clients {
		if client.id == event.TargetID {
			client.send <- &Action{ActionType: "ADD_INVITE", Payload: event}
		}
	}
}

func (h *WSHub) inviteResponded(event events.InviteRespondedEvent) {
	for client := range h.clients {
		if client.id == event.IssuerID {
			client.send <- &Action{ActionType: "UPDATE_INVITE", Payload: event}
		}
	}
}

func (h *WSHub) messageDeleted(event events.MessageDeletedEvent) {
	for client := range h.clients {
		if _, ok := client.groups[event.GroupID]; ok {
			client.send <- &Action{ActionType: "DELETE_MESSAGE", Payload: event}
		}
	}
}
