package eventprocessor

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
)

// EventProcessor processes traffic from listener and updates state of application
type RelayEventProcessor struct {
	Listener  msgqueue.EventListener
	RelayChan chan<- msgqueue.Event
}

// NewEventProcessor is a constructor for EventProcessor type
func NewRelayEventProcessor(listener msgqueue.EventListener, relayChan chan<- msgqueue.Event) *RelayEventProcessor {
	return &RelayEventProcessor{
		Listener:  listener,
		RelayChan: relayChan,
	}
}

// ProcessEvents listens to listener and updates state of application
func (p *RelayEventProcessor) ProcessEvents(eventNames ...string) {

	received, errors, err := p.Listener.Listen(eventNames...)
	if err != nil {
		log.Fatalf("Error when starting listening to kafka: %v", err)
	}

	for {
		select {
		case evt := <-received:
			p.RelayChan <- evt
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
