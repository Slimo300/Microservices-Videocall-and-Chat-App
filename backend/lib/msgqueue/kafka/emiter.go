package kafka

import (
	"strings"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
)

type kafkaEventEmiter struct {
	producer sarama.SyncProducer
	encoder  msgqueue.Encoder
}

type kafkaMessage struct {
	EventName string      `json:"eventName"`
	Payload   interface{} `json:"payload"`
}

// NewKafkaEventEmiter creates kafka emiter
func NewKafkaEventEmiter(client sarama.Client) (msgqueue.EventEmiter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &kafkaEventEmiter{
		producer: producer,
		encoder:  msgqueue.NewJSONEncoder(),
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

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageBody),
	})
	return err
}
