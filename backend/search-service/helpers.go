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

// kafkaSetup starts kafka listener
func kafkaSetup(brokerAddresses []string) (msgqueue.EventListener, error) {

	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "searchService"
	brokerConf.Version = sarama.V2_3_0_0
	client, err := sarama.NewClient(brokerAddresses, brokerConf)
	if err != nil {
		return nil, err
	}

	mapper := msgqueue.NewDynamicEventMapper()
	if err := mapper.RegisterTypes(
		reflect.TypeOf(events.UserRegisteredEvent{}),
		reflect.TypeOf(events.UserPictureModifiedEvent{}),
	); err != nil {
		return nil, err
	}

	listener, err := kafka.NewConsumerGroupEventListener(client, "search-service", mapper, &kafka.ListenerOptions{
		Logger: log.New(os.Stdout, "[listener]: ", log.Flags()),
	})
	if err != nil {
		return nil, err
	}

	return listener, nil
}
