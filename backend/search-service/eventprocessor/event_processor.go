package eventprocessor

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/database"
)

// EventProcessor listens to events coming from broker and updates state of application
type EventProcessor struct {
	DB       database.DBLayer
	Listener msgqueue.EventListener
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(db database.DBLayer, listener msgqueue.EventListener) *EventProcessor {
	return &EventProcessor{
		DB:       db,
		Listener: listener,
	}
}

func (p *EventProcessor) ProcessEvents(eventNames ...string) {

	received, errors, err := p.Listener.Listen(eventNames...)
	if err != nil {
		log.Fatalf("Error when starting listening to kafka: %v", err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.UserRegisteredEvent:
				if err := p.DB.AddUser(*e); err != nil {
					log.Printf("Adding user returned error: %v", err)
				}
			case *events.UserPictureModifiedEvent:
				if err := p.DB.UpdateProfilePicture(*e); err != nil {
					log.Printf("Updating profile picture url returned error: %v", err)
				}
			default:
				log.Println("Unsupported event type")
			}
		case err = <-errors:
			log.Printf("Error from listener: %v", err)
		}
	}
}
