package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
)

type kafkaEventListener struct {
	consumer sarama.Consumer
	mapper   msgqueue.EventMapper
	topics   []KafkaTopic
	Decoder  msgqueue.Decoder
	Offset   int64
}

type KafkaTopic struct {
	Name       string
	Partitions []int32
}

func NewKafkaEventListener(client sarama.Client, mapper msgqueue.EventMapper, topics ...KafkaTopic) (*kafkaEventListener, error) {

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &kafkaEventListener{
		consumer: consumer,
		mapper:   mapper,
		topics:   topics,
		Decoder:  msgqueue.NewJSONDecoder(),
		Offset:   sarama.OffsetNewest,
	}, nil

}

func (k *kafkaEventListener) Listen(events ...string) (<-chan msgqueue.Event, <-chan error, error) {

	var err error
	results := make(chan msgqueue.Event)
	errors := make(chan error)

	for _, topic := range k.topics {

		// When partitions is an empty slice the listener will listen on all of partitions
		partitions := topic.Partitions
		if len(partitions) == 0 {
			partitions, err = k.consumer.Partitions(topic.Name)
			if err != nil {
				return nil, nil, err
			}
		}

		for _, partition := range partitions {

			con, err := k.consumer.ConsumePartition(topic.Name, partition, k.Offset)
			if err != nil {
				return nil, nil, err
			}

			go func() {
				for msg := range con.Messages() {
					body := kafkaMessage{}
					err := k.Decoder.Decode(msg.Value, &body)
					if err != nil {
						errors <- fmt.Errorf("Could not unmarshal message: %s", err.Error())
						continue
					}
					evt, err := k.mapper.MapEvent(body.EventName, body.Payload)
					if err != nil {
						errors <- fmt.Errorf("Error when mapping events: %s", err.Error())
						continue
					}
					results <- evt
				}
			}()

			go func() {
				for err := range con.Errors() {
					errors <- err
				}
			}()
		}

	}

	return results, errors, err
}
