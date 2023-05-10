package eventprocessor

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/storage"
)

// EventProcessor listens to events coming from broker and updates state of application
type EventProcessor struct {
	DB       database.DBLayer
	Listener msgqueue.EventListener
	Storage  storage.StorageLayer
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(listener msgqueue.EventListener, db database.DBLayer, storage storage.StorageLayer) *EventProcessor {
	return &EventProcessor{
		DB:       db,
		Listener: listener,
		Storage:  storage,
	}
}

func (p *EventProcessor) ProcessEvents(eventNames ...string) {

	received, errors, err := p.Listener.Listen(eventNames...)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.MessageSentEvent:
				if err := p.DB.AddMessage(*e); err != nil {
					log.Printf("Error when adding message: %s\n", err.Error())
				}
			case *events.GroupDeletedEvent:
				if err := p.DB.DeleteGroupMembers(*e); err != nil {
					log.Printf("Error when deleting group members: %s\n", err.Error())
				}
				go func() {
					if err := p.Storage.DeleteFolder(e.ID.String()); err != nil {
						log.Printf("Error when deleting group files from storage: %s\n", err.Error())
					}
				}()
			case *events.MemberCreatedEvent:
				if err := p.DB.NewMember(*e); err != nil {
					log.Printf("Error when creating member from message: %s\n", err.Error())
				}
			case *events.MemberUpdatedEvent:
				if err := p.DB.ModifyMember(*e); err != nil {
					log.Printf("Error when updating member from message: %s\n", err.Error())
				}
			case *events.MemberDeletedEvent:
				if err := p.DB.DeleteMember(*e); err != nil {
					log.Printf("Error when deleting member from message: %s\n", err.Error())
				}
			default:
				log.Println("Event type not known")
			}
		case err = <-errors:
			log.Printf("Error when receiving message: %s", err.Error())
		}
	}

}
