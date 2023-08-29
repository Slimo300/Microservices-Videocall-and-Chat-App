package main

import (
	"log"
	"os"
	"reflect"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/kafka"
)

type ServiceStartedEvent struct {
	FQDN string `json:"fqdn" mapstructure:"fqdn"`
}

func (ServiceStartedEvent) EventName() string {
	return "webrtc.service-start"
}

// kafkaSetup starts kafka EventEmiter and EventListener
func kafkaSetup(brokerAddreses []string) (emiter msgqueue.EventEmiter, dbListener msgqueue.EventListener, relayListener msgqueue.EventListener, err error) {
	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "webrtcService"
	brokerConf.Version = sarama.V2_3_0_0
	brokerConf.Producer.Return.Successes = true
	client, err := sarama.NewClient(brokerAddreses, brokerConf)
	if err != nil {
		return nil, nil, nil, err
	}

	// initializing emiter
	emiter, err = kafka.NewKafkaEventEmiter(client, log.New(os.Stdout, "[ emiter ]: ", log.Flags()))
	if err != nil {
		return nil, nil, nil, err
	}

	// initializing dbListener
	dbListenerMapper := msgqueue.NewDynamicEventMapper()
	if err := dbListenerMapper.RegisterTypes(
		reflect.TypeOf(events.GroupDeletedEvent{}),
		reflect.TypeOf(events.MemberCreatedEvent{}),
		reflect.TypeOf(events.MemberDeletedEvent{}),
	); err != nil {
		return nil, nil, nil, err
	}

	dbListener, err = kafka.NewConsumerGroupEventListener(client, "webrtc-service", dbListenerMapper, &kafka.ListenerOptions{
		Logger: log.New(os.Stdout, "[DB listener]: ", log.Flags()),
	})

	// initializing relayListener
	relayListenerMapper := msgqueue.NewDynamicEventMapper()
	if err := relayListenerMapper.RegisterTypes(
		reflect.TypeOf(events.GroupDeletedEvent{}),
		reflect.TypeOf(events.MemberCreatedEvent{}),
		reflect.TypeOf(events.MemberDeletedEvent{}),
		reflect.TypeOf(events.MemberUpdatedEvent{}),
	); err != nil {
		return nil, nil, nil, err
	}

	relayListener, err = kafka.NewBroadcastEventListener(client, relayListenerMapper, &kafka.ListenerOptions{
		Logger: log.New(os.Stdout, "[Relay listener]: ", log.Flags()),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return emiter, dbListener, relayListener, nil
}
