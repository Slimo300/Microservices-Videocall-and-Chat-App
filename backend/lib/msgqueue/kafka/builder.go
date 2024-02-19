package kafka

import (
	"fmt"
	"reflect"

	"github.com/IBM/sarama"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
)

type KafkaBuilder struct {
	client sarama.Client
}

func NewKafkaBuilder(brokerAddresses []string) (msgqueue.BrokerBuilder, error) {

	brokerConf := sarama.NewConfig()
	brokerConf.Version = sarama.V2_3_0_0

	brokerConf.Producer.Return.Successes = true
	client, err := sarama.NewClient(brokerAddresses, brokerConf)
	if err != nil {
		return nil, err
	}

	return &KafkaBuilder{
		client: client,
	}, nil
}

func (b *KafkaBuilder) GetEmiter(conf msgqueue.EmiterConfig) (msgqueue.EventEmiter, error) {
	emiter, err := NewKafkaEventEmiter(b.client, conf.Logger)
	if err != nil {
		return nil, err
	}

	return emiter, nil
}

func (b *KafkaBuilder) GetListener(conf msgqueue.ListenerConfig) (listener msgqueue.EventListener, err error) {

	eventMapper := msgqueue.NewDynamicEventMapper()
	for _, ev := range conf.Events {
		if err := eventMapper.RegisterEventType(reflect.TypeOf(ev)); err != nil {
			return nil, fmt.Errorf("error registering event type: %w", err)
		}
	}

	if conf.Broadcast {
		listener, err = NewBroadcastEventListener(b.client, eventMapper, &ListenerOptions{
			Logger:        conf.Logger,
			SetPartitions: conf.SetPartitions,
			Decoder:       conf.Decoder,
		})
		if err != nil {
			return nil, fmt.Errorf("error creating broadcast listener: %w", err)
		}

	} else {
		listener, err = NewConsumerGroupEventListener(b.client, conf.ClientName, eventMapper, &ListenerOptions{
			Logger:        conf.Logger,
			SetPartitions: conf.SetPartitions,
			Decoder:       conf.Decoder,
		})
		if err != nil {
			return nil, fmt.Errorf("error creating consumer group listener: %w", err)
		}
	}

	return
}
