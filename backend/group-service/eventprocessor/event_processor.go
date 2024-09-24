package eventprocessor

import (
	"context"
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
)

// EventProcessor processes events from listener and updates state of application
type EventProcessor struct {
	DB       database.GroupsRepository
	Listener msgqueue.EventListener
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(db database.GroupsRepository, listener msgqueue.EventListener) *EventProcessor {
	return &EventProcessor{
		DB:       db,
		Listener: listener,
	}
}

// Process events listens to listener and updates state of application
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
				if _, err := p.DB.CreateUser(context.Background(), &models.User{ID: e.ID, UserName: e.Username}); err != nil {
					log.Printf("Listener NewUser error: %s", err.Error())
				}
			case *events.UserPictureModifiedEvent:
				if _, err := p.DB.UpdateUser(context.Background(), &models.User{ID: e.ID, Picture: e.PictureURL}); err != nil {
					log.Printf("Listener UpdatePicture error: %s", err.Error())
				}
			default:
				log.Println("Unsupported event type")
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
