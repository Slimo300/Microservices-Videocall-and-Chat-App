package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
)

type consumerGroupEventListener struct {
	client   sarama.Client
	consumer sarama.ConsumerGroup
	groupID  string
	mapper   msgqueue.EventMapper

	decoder msgqueue.Decoder
	offset  int64
}

func (c *consumerGroupEventListener) applyOptions(options *ListenerOptions) error {
	if options == nil {
		return nil
	}

	if options.Decoder != nil {
		c.decoder = *options.Decoder
	}

	if options.Offset != nil {
		c.offset = *options.Offset
	}

	if options.SetPartitions != nil {
		if err := c.setPartitions(options.SetPartitions); err != nil {
			return err
		}
	}

	return nil
}

func (c *consumerGroupEventListener) setPartitions(assignments map[string]int32) error {
	coor, err := c.client.Coordinator(c.groupID)
	if err != nil {
		return err
	}

	topicPartitions := make(map[string]*sarama.TopicPartition)

	for topic, num_of_partitions := range assignments {
		topicPartitions[topic] = &sarama.TopicPartition{
			Count: num_of_partitions,
		}
	}

	if _, err := coor.CreatePartitions(&sarama.CreatePartitionsRequest{
		TopicPartitions: topicPartitions,
		Timeout:         10 * time.Second,
		ValidateOnly:    false,
	}); err != nil {
		return err
	}

	return nil
}

// Returns new consumerGroupEventListener
func NewConsumerGroupEventListener(client sarama.Client, groupID string, mapper msgqueue.EventMapper, options *ListenerOptions) (msgqueue.EventListener, error) {

	listener := &consumerGroupEventListener{
		client:  client,
		groupID: groupID,
		mapper:  mapper,
		decoder: msgqueue.NewJSONDecoder(),
		offset:  OffsetNewest,
	}

	if err := listener.applyOptions(options); err != nil {
		return nil, err
	}

	return listener, nil
}

// Listen listens for specified topics and sends them through channel
func (c *consumerGroupEventListener) Listen(topics ...string) (<-chan msgqueue.Event, <-chan error, error) {

	results := make(chan msgqueue.Event)
	errors := make(chan error)

	consumer, err := sarama.NewConsumerGroupFromClient(c.groupID, c.client)
	if err != nil {
		return nil, nil, err
	}

	groupConsumer := &groupConsumer{
		ready:   make(chan bool),
		results: make(chan *sarama.ConsumerMessage),
	}

	go func() {
		for {
			if err := consumer.Consume(context.TODO(), topics, groupConsumer); err != nil {
				return
			}
			groupConsumer.ready = make(chan bool)
		}
	}()

	<-groupConsumer.ready

	go func() {
		for msg := range groupConsumer.results {

			body := kafkaMessage{}
			if err := c.decoder.Decode(msg.Value, &body); err != nil {
				errors <- fmt.Errorf("Could not unmarshal message: %s", err.Error())
				continue
			}

			evt, err := c.mapper.MapEvent(body.EventName, body.Payload)
			if err != nil {
				errors <- fmt.Errorf("Error when mapping event: %s", err.Error())
				continue
			}

			results <- evt
		}
	}()

	return results, errors, nil
}

// Group consumer
type groupConsumer struct {
	ready   chan bool
	results chan *sarama.ConsumerMessage
}

func (consumer *groupConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *groupConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	close(consumer.results)
	return nil
}

func (consumer *groupConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			consumer.results <- message
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
