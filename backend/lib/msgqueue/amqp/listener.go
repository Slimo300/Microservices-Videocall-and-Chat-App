package amqp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue      string
	mapper     msgqueue.EventMapper
	decoder    msgqueue.Decoder
}

// NewAMQPEventListener creates amqp listener
func NewAMQPEventListener(conn *amqp.Connection, mapper msgqueue.EventMapper, queueName string, exchanges ...string) (msgqueue.EventListener, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	for _, ex := range exchanges {
		if err := channel.ExchangeDeclare(ex, "topic", true, false, false, false, nil); err != nil {
			return nil, err
		}
	}

	return &amqpEventListener{
		connection: conn,
		queue:      queueName,
		mapper:     mapper,
		decoder:    msgqueue.NewJSONDecoder(),
	}, nil
}

// Listen will check amqp queue for given events and send them through returned channel
func (a *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	eventChan := make(chan msgqueue.Event)
	errChan := make(chan error)

	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	// defer channel.Close()

	// Here we bind listener queue to exchanges via routing keys provided in 'eventNames' argument,
	// event is routing key and its first part is name of exchange it will be published to e.g.:
	// event - users.created -> exchange users
	for _, event := range eventNames {
		channel.QueueBind(a.queue, event, strings.Split(event, ".")[0], false, nil)
	}

	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for msg := range msgs {
			evtName, ok := msg.Headers["x-event-name"]
			if !ok {
				errChan <- errors.New("message did not contain x-event-name header")
				msg.Nack(false, false)
				continue
			}

			eventName, ok := evtName.(string)
			if !ok {
				errChan <- fmt.Errorf("header %s did not contain string", eventName)
				msg.Nack(false, false)
				continue
			}

			var messageBody interface{}
			err := a.decoder.Decode(msg.Body, messageBody)
			if err != nil {
				errChan <- fmt.Errorf("decoding message returned error: %v", err)
				msg.Nack(false, false)
				continue
			}

			event, err := a.mapper.MapEvent(eventName, messageBody)
			if err != nil {
				errChan <- fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
				msg.Nack(false, false)
				continue
			}

			eventChan <- event
			msg.Ack(false)

		}
	}()

	return eventChan, errChan, nil
}
