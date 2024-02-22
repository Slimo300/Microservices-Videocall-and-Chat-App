package amqp

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type amqpEventEmiter struct {
	connection *amqp.Connection
	exchange   string
	encoder    msgqueue.Encoder
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
		exchange:   exchange,
		encoder:    msgqueue.NewJSONEncoder(),
	}, nil
}

// Emit sends a new Event to amqp
func (e *amqpEventEmiter) Emit(evt msgqueue.Event) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	jsonBody, err := e.encoder.Encode(evt)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": evt.EventName()},
		ContentType: "application/json",
		Body:        jsonBody,
	}

	if err := channel.PublishWithContext(context.Background(), e.exchange, evt.EventName(), false, false, msg); err != nil {
		return err
	}

	return nil
}
