package amqp

import (
	"encoding/json"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/streadway/amqp"
)

type amqpEventEmitter struct {
	connection *amqp.Connection
	exchange   string
}

func NewAMQPEventEmitter(conn amqp.Connection, exchange string) (*amqpEventEmitter, error) {

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer channel.Close()

	if err := channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}
	return &amqpEventEmitter{
		connection: &conn,
		exchange:   exchange,
	}, nil
}

func (a *amqpEventEmitter) Emit(evt msgqueue.Event) error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	jsonBody, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": evt.EventName()},
		ContentType: "application/json",
		Body:        jsonBody,
	}

	if err := channel.Publish(a.exchange, evt.EventName(), false, false, msg); err != nil {
		return err
	}

	return nil
}
