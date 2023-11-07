package ws

import (
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
)

func (h *WSHub) HandleEvent(event msgqueue.Event) {
	switch e := event.(type) {
	case *events.GroupDeletedEvent:
		h.groupDeleted(*e)
	case *events.MemberCreatedEvent:
		h.memberAdded(*e)
	case *events.MemberDeletedEvent:
		h.memberDeleted(*e)
	case *events.MemberUpdatedEvent:
		h.memberUpdated(*e)
	case *events.InviteSentEvent:
		h.inviteSent(*e)
	case *events.InviteRespondedEvent:
		h.inviteResponded(*e)
	case *events.MessageDeletedEvent:
		h.messageDeleted(*e)
	case *events.MessageSentEvent:
		h.messageSent(*e)
	default:
		log.Println("Unsupported Event Type: ", event.EventName())
	}
}

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
				ID:       event.ID,
				MemberID: event.MemberID,
				Member: Member{
					ID:       event.GroupID,
					GroupID:  event.GroupID,
					UserID:   event.UserID,
					Username: event.Nick,
				},
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
