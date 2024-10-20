package eventprocessor

import (
	"context"
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app/command"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
)

// EventProcessor processes events from listener and updates state of application
type EventProcessor struct {
	app      app.App
	listener msgqueue.EventListener
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(app app.App, listener msgqueue.EventListener) *EventProcessor {
	return &EventProcessor{
		app:      app,
		listener: listener,
	}
}

// Process events listens to listener and updates state of application
func (p *EventProcessor) ProcessEvents(eventNames ...string) {
	received, errors, err := p.listener.Listen(eventNames...)
	if err != nil {
		log.Fatalf("Error when starting listening to kafka: %v", err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.UserVerifiedEvent:
				if err := p.app.Commands.CreateUser.Handle(context.Background(), command.CreateUserCommand{
					UserID:   e.ID,
					Username: e.Username,
				}); err != nil {
					log.Printf("Listener error: %v", err)
				}
			case *events.UserPictureModifiedEvent:
				if err := p.app.Commands.UpdateUser.Handle(context.Background(), command.UpdateUserCommand{
					UserID:     e.ID,
					HasPicture: e.HasPicture,
				}); err != nil {
					log.Printf("Listener error: %v", err)
				}
			default:
				log.Println("Unsupported event type")
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
