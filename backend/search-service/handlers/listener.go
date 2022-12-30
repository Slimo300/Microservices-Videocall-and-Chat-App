package handlers

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

func (s *Server) RunListener() {

	eventChan, errorChan, err := s.Listener.Listen()
	if err != nil {
		log.Fatalf("Listener couldn't launch: %v", err)
	}

	for {
		select {
		case evt := <-eventChan:
			switch e := evt.(type) {
			case *events.UserRegisteredEvent:
				if err := s.DB.AddUser(*e); err != nil {
					log.Printf("Adding user returned error: %v", err)
				}
			}
		case err = <-errorChan:
			log.Printf("Error from listener: %v", err)
		}
	}
}
