package amqp

import (
	"context"
	"strings"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type amqpEventEmiter struct {
	connection *amqp.Connection
	Encoder    msgqueue.Encoder
}

// NewAMQPEventEmiter creates amqp emiter
func NewAMQPEventEmiter(conn *amqp.Connection, exchange string) (msgqueue.EventEmiter, error) {

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer channel.Close()

	if err := channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}
	return &amqpEventEmiter{
		connection: conn,
		Encoder:    msgqueue.NewJSONEncoder(),
	}, nil
}

// Emit sends a new Event to amqp
func (a *amqpEventEmiter) Emit(evt msgqueue.Event) error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	jsonBody, err := a.Encoder.Encode(evt)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": evt.EventName()},
		ContentType: "application/json",
		Body:        jsonBody,
	}

	if err := channel.PublishWithContext(context.Background(), strings.Split(evt.EventName(), ".")[0], evt.EventName(), false, false, msg); err != nil {
		return err
	}

	return nil
}
