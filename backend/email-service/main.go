package main

import (
	"log"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/email-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/email-service/eventprocessor"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/email-service/service"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/amqp"
)

func main() {
	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Error when reading configuration: %v", err)
	}

	emailService, err := service.NewEmailService(conf.EmailFrom,
		conf.SMTPHost,
		conf.SMTPPort,
		conf.SMTPUser,
		conf.SMTPPass,
		conf.Origin,
	)
	if err != nil {
		log.Fatalf("Error when creating email service: %v", err)
	}

	builder, err := amqp.NewAMQPBuilder(conf.BrokerAddress)
	if err != nil {
		log.Fatalf("Error when creating amqp builder: %v", err)
	}
	listener, err := builder.GetListener(msgqueue.ListenerConfig{
		ClientName: "email-service",
		Events: []msgqueue.Event{
			events.UserRegisteredEvent{},
			events.UserForgotPasswordEvent{},
		},
	})
	if err != nil {
		log.Fatalf("Error when creating amqp listener: %v", err)
	}

	eventprocessor.NewEventProcessor(listener, emailService).ProcessEvents("user")
}
