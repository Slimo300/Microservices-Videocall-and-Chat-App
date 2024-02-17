package main

import (
	"log"
	"os"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/kafka"

	"github.com/IBM/sarama"
)

// kafkaSetup starts Kafka EventEmiter and EventListener
func kafkaSetup(brokerAddresses []string) (msgqueue.EventEmiter, error) {

	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "usersService"
	brokerConf.Version = sarama.V2_3_0_0
	brokerConf.Producer.Return.Successes = true
	client, err := sarama.NewClient(brokerAddresses, brokerConf)
	if err != nil {
		return nil, err
	}

	emiter, err := kafka.NewKafkaEventEmiter(client, log.New(os.Stdout, "[ emiter ]: ", log.Flags()))
	if err != nil {
		return nil, err
	}

	return emiter, nil

}
