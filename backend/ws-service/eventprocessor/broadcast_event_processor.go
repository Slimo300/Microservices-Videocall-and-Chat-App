package eventprocessor

import (
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/database"
)

// EventProcessor processes traffic from listener and updates state of application
type HubEventProcessor struct {
	Listener msgqueue.EventListener
	DB       database.DBLayer
	HubChan  chan<- msgqueue.Event
}

// NewEventProcessor is a constructor for EventProcessor type
func NewHubEventProcessor(listener msgqueue.EventListener, hubchan chan<- msgqueue.Event) *HubEventProcessor {
	return &HubEventProcessor{
		Listener: listener,
		HubChan:  hubchan,
	}
}

// ProcessEvents listens to listener and updates state of application
func (p *HubEventProcessor) ProcessEvents(eventNames ...string) {

	received, errors, err := p.Listener.Listen(eventNames...)
	if err != nil {
		log.Fatalf("Error when starting listening to kafka: %v", err)
	}

	for {
		select {
		case evt := <-received:
			p.HubChan <- evt
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
