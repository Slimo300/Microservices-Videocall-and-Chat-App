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

// kafkaSetup starts kafka EventEmiter and EventListener
func kafkaSetup(brokerAddreses []string) (msgqueue.EventListener, error) {
	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "webrtcGatewayService"
	brokerConf.Version = sarama.V2_3_0_0
	brokerConf.Producer.Return.Successes = true
	client, err := sarama.NewClient(brokerAddreses, brokerConf)
	if err != nil {
		return nil, err
	}

	// initializing dbListener
	mapper := msgqueue.NewDynamicEventMapper()
	if err := mapper.RegisterTypes(
		reflect.TypeOf(events.ServiceStartedEvent{}),
	); err != nil {
		return nil, err
	}

	listener, err := kafka.NewConsumerGroupEventListener(client, "webrtc-gateway-service", mapper, &kafka.ListenerOptions{
		Logger: log.New(os.Stdout, "[listener]: ", log.Flags()),
	})
	if err != nil {
		return nil, err
	}

	return listener, nil
}
