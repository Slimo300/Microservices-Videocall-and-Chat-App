package handlers

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/events"
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
			case *events.UserRegisteredEvent:
				if err := s.DB.NewUser(*e); err != nil {
					log.Printf("Listener NewUser error: %s", err.Error())
				}
			default:
				log.Println("Unsupported event type")
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
