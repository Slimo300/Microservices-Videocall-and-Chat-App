package kafka

import (
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
)

type kafkaEventEmiter struct {
	producer sarama.SyncProducer
	encoder  msgqueue.Encoder
	logger   *log.Logger
}

type kafkaMessage struct {
	EventName string      `json:"eventName"`
	Payload   interface{} `json:"payload"`
}

// NewKafkaEventEmiter creates kafka emiter
func NewKafkaEventEmiter(client sarama.Client, logger *log.Logger) (msgqueue.EventEmiter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &kafkaEventEmiter{
		producer: producer,
		encoder:  msgqueue.NewJSONEncoder(),
		logger:   logger,
	}, nil
}

// Emit sends a new Event to kafka
func (k *kafkaEventEmiter) Emit(event msgqueue.Event) error {
	messageBody, err := k.encoder.Encode(kafkaMessage{
		EventName: event.EventName(),
		Payload:   event,
	})
	if err != nil {
		return err
	}

	topic := strings.Split(event.EventName(), ".")[0]

	if k.logger != nil {
		k.logger.Printf("Emiting: %s\n", event.EventName())
	}

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageBody),
	})
	return err
}
