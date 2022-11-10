package handlers

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

func (s *Server) RunListener(eventNames ...string) {

	received, errors, err := s.Listener.Listen(eventNames...)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.GroupDeletedEvent:
				if err := s.DB.DeleteGroupMembers(*e); err != nil {
					log.Printf("Listener DeleteGroup error: %s", err.Error())
				}
				s.Hub.GroupDeleted(*e)
			case *events.MemberCreatedEvent:
				if err := s.DB.NewMember(*e); err != nil {
					log.Printf("Listener NewMember error: %s", err.Error())
				}
				s.Hub.MemberAdded(*e)
			case *events.MemberDeletedEvent:
				if err := s.DB.DeleteMember(*e); err != nil {
					log.Printf("Listener DeleteMember error: %s", err.Error())
				}
				s.Hub.MemberDeleted(*e)
			case *events.MemberUpdatedEvent:
				s.Hub.MemberUpdated(*e)
			case *events.InviteSentEvent:
				s.Hub.InviteSent(*e)
			case *events.MessageDeletedEvent:
				s.Hub.MessageDeleted(*e)
			default:
				log.Println("Unsupported event type")
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
