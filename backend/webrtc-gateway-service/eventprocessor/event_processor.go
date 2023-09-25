package eventprocessor

import (
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-gateway-service/database"
)

// EventProcessor processes traffic from listener and updates state of application
type EventProcessor struct {
	Listener msgqueue.EventListener
	DB       database.DBLayer
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(listener msgqueue.EventListener, db database.DBLayer) *EventProcessor {
	return &EventProcessor{
		Listener: listener,
		DB:       db,
	}
}

// ProcessEvents listens to listener and updates state of application
func (p *EventProcessor) ProcessEvents(eventNames ...string) {

	received, errors, err := p.Listener.Listen(eventNames...)
	if err != nil {
		log.Fatalf("Error when starting listening to kafka: %v", err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.ServiceStartedEvent:
				log.Printf("New Service Instance: %v", e.ServiceAddress)
				if err := p.DB.NewInstance(e.ServiceAddress); err != nil {
					log.Printf("Error adding new instance to database: %v", err)
				}
			default:
				log.Printf("Unsupported event type")
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
