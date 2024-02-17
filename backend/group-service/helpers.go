package main

import (
	"log"
	"os"
	"reflect"

	"github.com/IBM/sarama"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/kafka"
)

// kafkaSetup starts Kafka EventEmiter and EventListener
func kafkaSetup(brokerAddresses []string) (msgqueue.EventEmiter, msgqueue.EventListener, error) {

	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "groupsService"
	brokerConf.Version = sarama.V2_3_0_0
	brokerConf.Producer.Return.Successes = true
	client, err := sarama.NewClient(brokerAddresses, brokerConf)
	if err != nil {
		return nil, nil, err
	}

	emiter, err := kafka.NewKafkaEventEmiter(client, log.New(os.Stdout, "[ emiter ]: ", log.Flags()))
	if err != nil {
		return nil, nil, err
	}
	mapper := msgqueue.NewDynamicEventMapper()
	if err := mapper.RegisterTypes(
		reflect.TypeOf(events.UserRegisteredEvent{}),
		reflect.TypeOf(events.UserPictureModifiedEvent{}),
	); err != nil {
		return nil, nil, err
	}
	listener, err := kafka.NewConsumerGroupEventListener(client, "group-service", mapper, &kafka.ListenerOptions{
		Logger: log.New(os.Stdout, "[listener]: ", log.Flags()),
	})
	if err != nil {
		return nil, nil, err
	}

	return emiter, listener, nil

}
