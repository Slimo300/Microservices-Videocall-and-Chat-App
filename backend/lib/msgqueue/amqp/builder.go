package amqp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPBuilder struct {
	conn *amqp.Connection
}

func NewAMQPBuilder(brokerAddress string) (msgqueue.BrokerBuilder, error) {

	conn, err := amqp.Dial(brokerAddress)
	if err != nil {
		return nil, err
	}

	return &AMQPBuilder{
		conn: conn,
	}, nil
}

func (b *AMQPBuilder) GetEmiter(conf msgqueue.EmiterConfig) (msgqueue.EventEmiter, error) {
	emiter, err := NewAMQPEventEmiter(b.conn, conf.ExchangeName)
	if err != nil {
		return nil, err
	}

	return emiter, nil
}

func (b *AMQPBuilder) GetListener(conf msgqueue.ListenerConfig) (listener msgqueue.EventListener, err error) {

	exchanges := make(map[string]bool)

	eventMapper := msgqueue.NewDynamicEventMapper()
	for _, ev := range conf.Events {
		if err := eventMapper.RegisterEventType(reflect.TypeOf(ev)); err != nil {
			return nil, fmt.Errorf("error registering event type: %w", err)
		}

		exchanges[strings.Split(ev.EventName(), ".")[0]] = true
	}

	var exchangesList []string
	for ex := range exchanges {
		exchangesList = append(exchangesList, ex)
	}

	listener, err = NewAMQPEventListener(b.conn, eventMapper, conf.ClientName, exchangesList...)
	if err != nil {
		return nil, err
	}

	return listener, err
}
