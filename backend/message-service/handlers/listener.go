package handlers

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/events"
)

func (s *Server) RunListener(eventNames ...string) {

	received, errors, err := s.Listener.Listen(eventNames...)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.MessageSentEvent:
				if err := s.DB.AddMessage(*e); err != nil {
					log.Printf("Error when adding message: %s\n", err.Error())
				}
			case *events.GroupDeletedEvent:
				if err := s.DB.DeleteGroupMembers(*e); err != nil {
					log.Printf("Error when deleting group members: %s\n", err.Error())
				}
			case *events.MemberCreatedEvent:
				if err := s.DB.NewMember(*e); err != nil {
					log.Printf("Error when creating member from message: %s\n", err.Error())
				}
			case *events.MemberUpdatedEvent:
				if err := s.DB.ModifyMember(*e); err != nil {
					log.Printf("Error when updating member from message: %s\n", err.Error())
				}
			case *events.MemberDeletedEvent:
				if err := s.DB.DeleteMember(*e); err != nil {
					log.Printf("Error when deleting member from message: %s\n", err.Error())
				}
			default:
				log.Println("Event type not known")
			}
		case err = <-errors:
			log.Printf("Error when receiving message: %s", err.Error())
		}
	}
}
