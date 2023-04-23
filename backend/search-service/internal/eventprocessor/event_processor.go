package eventprocessor

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/chat-searchservice/internal/database"
)

// EventProcessor listens to events coming from broker and updates state of application
type EventProcessor struct {
	DB       database.DBLayer
	Listener msgqueue.EventListener
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(listener msgqueue.EventListener, db database.DBLayer) *EventProcessor {
	return &EventProcessor{
		DB:       db,
		Listener: listener,
	}
}

func (p *EventProcessor) ProcessEvents(eventNames ...string) {

	eventChan, errorChan, err := p.Listener.Listen()
	if err != nil {
		log.Fatalf("Listener couldn't launch: %v", err)
	}

	for {
		select {
		case evt := <-eventChan:
			switch e := evt.(type) {
			case *events.UserRegisteredEvent:
				if err := p.DB.AddUser(*e); err != nil {
					log.Printf("Adding user returned error: %v", err)
				}
			case *events.UserPictureModifiedEvent:
				if err := p.DB.UpdateProfilePicture(*e); err != nil {
					log.Printf("Updating profile picture url returned error: %v", err)
				}
			}
		case err = <-errorChan:
			log.Printf("Error from listener: %v", err)
		}
	}
}
