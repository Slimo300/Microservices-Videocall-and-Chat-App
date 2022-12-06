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
				s.EventChan <- evt
			case *events.MemberCreatedEvent:
				if err := s.DB.NewMember(*e); err != nil {
					log.Printf("Listener NewMember error: %s", err.Error())
				}
				s.EventChan <- evt
			case *events.MemberDeletedEvent:
				if err := s.DB.DeleteMember(*e); err != nil {
					log.Printf("Listener DeleteMember error: %s", err.Error())
				}
				s.EventChan <- evt
			default:
				s.EventChan <- evt
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
