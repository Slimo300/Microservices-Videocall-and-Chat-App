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

// kafkaSetup starts kafka EventEmiter and EventListener
func kafkaSetup(brokerAddreses []string) (emiter msgqueue.EventEmiter, dbListener msgqueue.EventListener, hubListener msgqueue.EventListener, err error) {
	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "websocketService"
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

	dbListener, err = kafka.NewConsumerGroupEventListener(client, "ws-service", dbListenerMapper, &kafka.ListenerOptions{
		Logger: log.New(os.Stdout, "[DB listener]: ", log.Flags()),
	})

	// initializing hubListener
	hubListenerMapper := msgqueue.NewDynamicEventMapper()
	if err := hubListenerMapper.RegisterTypes(
		reflect.TypeOf(events.GroupDeletedEvent{}),
		reflect.TypeOf(events.MemberCreatedEvent{}),
		reflect.TypeOf(events.MemberDeletedEvent{}),
		reflect.TypeOf(events.MemberUpdatedEvent{}),
		reflect.TypeOf(events.MessageDeletedEvent{}),
		reflect.TypeOf(events.InviteSentEvent{}),
		reflect.TypeOf(events.InviteRespondedEvent{}),
		reflect.TypeOf(events.MessageSentEvent{}),
	); err != nil {
		return nil, nil, nil, err
	}

	hubListener, err = kafka.NewBroadcastEventListener(client, hubListenerMapper, &kafka.ListenerOptions{
		Logger: log.New(os.Stdout, "[Hub listener]: ", log.Flags()),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return emiter, dbListener, hubListener, nil
}
