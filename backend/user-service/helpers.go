package main

import (
	"log"
	"os"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	libamqp "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/amqp"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/kafka"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/IBM/sarama"
)

func rabbitSetup(brokerAddress string, exchangeName string) (msgqueue.EventEmiter, error) {
	conn, err := amqp.Dial(brokerAddress)
	if err != nil {
		return nil, err
	}

	emiter, err := libamqp.NewAMQPEventEmiter(conn, exchangeName)
	if err != nil {
		return nil, err
	}

	return emiter, nil
}

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
