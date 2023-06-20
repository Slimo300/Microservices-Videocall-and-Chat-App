package eventprocessor

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/database"
)

// EventProcessor processes traffic from listener and updates state of application
type DBEventProcessor struct {
	Listener msgqueue.EventListener
	DB       database.DBLayer
}

// NewEventProcessor is a constructor for EventProcessor type
func NewDBEventProcessor(listener msgqueue.EventListener, db database.DBLayer) *DBEventProcessor {
	return &DBEventProcessor{
		Listener: listener,
		DB:       db,
	}
}

// ProcessEvents listens to listener and updates state of application
func (p *DBEventProcessor) ProcessEvents(eventNames ...string) {

	received, errors, err := p.Listener.Listen(eventNames...)
	if err != nil {
		log.Fatalf("Error when starting listening to kafka: %v", err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.GroupDeletedEvent:
				if err := p.DB.DeleteGroup(*e); err != nil {
					log.Printf("Listener DeleteGroup error: %s", err.Error())
				}
			case *events.MemberCreatedEvent:
				if err := p.DB.NewMember(*e); err != nil {
					log.Printf("Listener NewMember error: %s", err.Error())
				}
			case *events.MemberDeletedEvent:
				if err := p.DB.DeleteMember(*e); err != nil {
					log.Printf("Listener DeleteMember error: %s", err.Error())
				}
			default:
				log.Printf("Unsupported event type")
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
