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

	return listener, nil
}
