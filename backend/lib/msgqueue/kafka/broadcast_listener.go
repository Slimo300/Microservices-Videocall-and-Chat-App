package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
)

const OffsetNewest = sarama.OffsetNewest
const OffsetOldest = sarama.OffsetOldest

// ListenerOptions allows to change default configuration of listener
type ListenerOptions struct {
	Decoder       *msgqueue.Decoder // default: msgqueue.jsonDecoder
	Offset        *int64            // default: sarama.OffsetNewest
	SetPartitions map[string]int32  // allows to set number of partitions for a topic, maps 'topic' -> no. of partitions
}

type broadcastEventListener struct {
	consumer sarama.Consumer
	mapper   msgqueue.EventMapper

	decoder msgqueue.Decoder
	offset  int64
}

func (b *broadcastEventListener) applyOptions(options *ListenerOptions) {
	if options == nil {
		return
	}

	if options.Decoder != nil {
		b.decoder = *options.Decoder
	}

	if options.Offset != nil {
		b.offset = *options.Offset
	}
}

// Returns new broadcastEventListener
func NewBroadcastEventListener(client sarama.Client, mapper msgqueue.EventMapper, options *ListenerOptions) (msgqueue.EventListener, error) {

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &broadcastEventListener{
		consumer: consumer,
		mapper:   mapper,
		decoder:  msgqueue.NewJSONDecoder(),
		offset:   OffsetNewest,
	}
	listener.applyOptions(options)

	return listener, nil
}

// Listen listens for specified topics and sends them through channel
func (b *broadcastEventListener) Listen(topics ...string) (<-chan msgqueue.Event, <-chan error, error) {

	results := make(chan msgqueue.Event)
	errors := make(chan error)

	for _, topic := range topics {
		partitions, err := b.consumer.Partitions(topic)
		if err != nil {
			return nil, nil, err
		}

		for _, partition := range partitions {
			con, err := b.consumer.ConsumePartition(topic, partition, b.offset)
			if err != nil {
				return nil, nil, err
			}

			go func() {
				for msg := range con.Messages() {
					body := kafkaMessage{}
					if err := b.decoder.Decode(msg.Value, &body); err != nil {
						errors <- fmt.Errorf("Could not unmarshal message: %s", err.Error())
						continue
					}

					evt, err := b.mapper.MapEvent(body.EventName, body.Payload)
					if err != nil {
						errors <- fmt.Errorf("Error when mapping event: %s", err.Error())
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

	return results, errors, nil
}
