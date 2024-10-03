package eventprocessor

import (
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/email-service/service"
)

// EventProcessor processes events from listener and updates state of application
type EventProcessor struct {
	Listener     msgqueue.EventListener
	EmailService service.EmailService
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(listener msgqueue.EventListener, emailService service.EmailService) *EventProcessor {
	return &EventProcessor{
		EmailService: emailService,
		Listener:     listener,
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
				if err := p.EmailService.SendVerificationEmail(e.Email, e.Username, e.Code); err != nil {
					log.Printf("Couldn't send verification error: %v", err)
				}
			case *events.UserForgotPasswordEvent:
				if err := p.EmailService.SendResetPasswordEmail(e.Email, e.Username, e.Code); err != nil {
					log.Printf("Couldn't send verification error: %v", err)
				}
			default:
				log.Println("Unsupported event type")
			}
		case err = <-errors:
			log.Printf("Listener error: %s", err.Error())
		}
	}
}
