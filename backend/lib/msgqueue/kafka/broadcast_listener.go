package kafka

import (
	"errors"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
)

const OffsetNewest = sarama.OffsetNewest
const OffsetOldest = sarama.OffsetOldest

// ListenerOptions allows to change default configuration of listener
type ListenerOptions struct {
	// default: msgqueue.jsonDecoder
	Decoder *msgqueue.Decoder
	// default: OffsetNewest (-1), 0 leaves default option, to start from first message written use OffsetOldest or -2
	Offset int64
	// allows to set number of partitions for a topic, maps 'topic' -> no. of partitions
	SetPartitions map[string]int32
	// allows to log to specific source
	Logger *log.Logger
}

type broadcastEventListener struct {
	consumer sarama.Consumer
	mapper   msgqueue.EventMapper

	decoder msgqueue.Decoder
	offset  int64
	logger  *log.Logger
}

func (b *broadcastEventListener) applyOptions(options *ListenerOptions) {
	if options == nil {
		return
	}

	if options.Decoder != nil {
		b.decoder = *options.Decoder
	}

	if options.Offset != 0 {
		b.offset = options.Offset
	}

	if options.Logger != nil {
		b.logger = options.Logger
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

	if len(topics) == 0 {
		return nil, nil, errors.New("Listen called with no topics provided")
	}

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

					if b.logger != nil {
						b.logger.Printf("Received %s\n", body.EventName)
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
