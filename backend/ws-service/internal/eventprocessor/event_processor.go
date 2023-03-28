package eventprocessor

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/chat-wsservice/internal/database"
)

// EventProcessor processes traffic from listener and updates state of application
type EventProcessor struct {
	Listener msgqueue.EventListener
	DB       database.DBLayer
	HubChan  chan<- msgqueue.Event
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(db database.DBLayer, listener msgqueue.EventListener, hubchan chan<- msgqueue.Event) *EventProcessor {
	return &EventProcessor{
		Listener: listener,
		DB:       db,
		HubChan:  hubchan,
	}
}

// ProcessEvents listens to listener and updates state of application
func (p *EventProcessor) ProcessEvents(eventNames ...string) {

	received, errors, err := p.Listener.Listen(eventNames...)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.GroupDeletedEvent:
				if err := p.DB.DeleteGroupMembers(*e); err != nil {
					log.Printf("Listener DeleteGroup error: %s", err.Error())
				}
				p.HubChan <- evt
			case *events.MemberCreatedEvent:
				if err := p.DB.NewMember(*e); err != nil {
					log.Printf("Listener NewMember error: %s", err.Error())
				}
				p.HubChan <- evt
			case *events.MemberDeletedEvent:
				if err := p.DB.DeleteMember(*e); err != nil {
					log.Printf("Listener DeleteMember error: %s", err.Error())
				}
				p.HubChan <- evt
			default:
				p.HubChan <- evt
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
