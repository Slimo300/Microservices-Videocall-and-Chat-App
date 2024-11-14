package eventprocessor

import (
	"context"
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/app"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/app/command"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
)

// EventProcessor listens to events coming from broker and updates state of application
type EventProcessor struct {
	listener msgqueue.EventListener
	app      app.App
}

// NewEventProcessor is a constructor for EventProcessor type
func NewEventProcessor(listener msgqueue.EventListener, app app.App) *EventProcessor {
	return &EventProcessor{
		app:      app,
		listener: listener,
	}
}

func (p *EventProcessor) ProcessEvents(eventNames ...string) {
	received, errors, err := p.listener.Listen(eventNames...)
	if err != nil {
		log.Fatalf("Error when starting listening to kafka: %v", err)
	}

	for {
		select {
		case evt := <-received:
			switch e := evt.(type) {
			case *events.MessageSentEvent:
				var files []models.MessageFile
				for _, file := range e.Files {
					files = append(files, models.NewMessageFile(e.ID, file.Key, file.Extension))
				}
				if err := p.app.Commands.CreateMessage.Handle(context.Background(), command.CreateMessageCommand{
					MessageID: e.ID,
					MemberID:  e.MemberID,
					GroupID:   e.GroupID,
					Text:      e.Text,
					Nick:      e.Nick,
					Posted:    e.Posted,
					Files:     files,
				}); err != nil {
					log.Printf("Error when adding message: %s\n", err.Error())
				}
			case *events.GroupDeletedEvent:
				if err := p.app.Commands.DeleteGroup.Handle(context.Background(), command.DeleteGroupCommand{GroupID: e.ID}); err != nil {
					log.Printf("Error when deleting group members: %s\n", err.Error())
				}
			case *events.MemberCreatedEvent:
				if err := p.app.Commands.CreateMember.Handle(context.Background(), command.CreateMemberCommand{
					MemberID: e.ID,
					GroupID:  e.GroupID,
					UserID:   e.UserID,
					Username: e.User.UserName,
					Creator:  e.Creator,
				}); err != nil {
					log.Printf("Error when creating member from event: %s\n", err.Error())
				}
			case *events.MemberUpdatedEvent:
				if err := p.app.Commands.UpdateMember.Handle(context.Background(), command.UpdateMemberCommand{
					MemberID:         e.ID,
					Admin:            e.Admin,
					DeletingMessages: e.DeletingMessages,
				}); err != nil {
					log.Printf("Error when updating member from event: %s\n", err.Error())
				}
			case *events.MemberDeletedEvent:
				if err := p.app.Commands.DeleteMember.Handle(context.Background(), command.DeleteMemberCommand{MemberID: e.ID}); err != nil {
					log.Printf("Error when deleting member from event: %s\n", err.Error())
				}
			default:
				log.Printf("Unknown event type: %s\n", e.EventName())
			}
		case err = <-errors:
			log.Printf("Error when receiving message: %s\n", err.Error())
		}
	}

}
