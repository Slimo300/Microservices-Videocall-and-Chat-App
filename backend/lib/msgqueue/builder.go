package msgqueue

import (
	"log"
)

type BrokerType int

const (
	KAFKA_BROKER = iota + 1
	AMQP_BROKER
)

// Possible types are "AMQP" or "KAFKA", otherwise function panics
func ParseBrokerType(brokerType string) BrokerType {
	switch brokerType {
	case "KAFKA":
		return KAFKA_BROKER
	case "AMQP":
		return AMQP_BROKER
	default:
		panic("Unsupported broker type")
	}
}

type BrokerBuilder interface {
	GetEmiter(EmiterConfig) (EventEmiter, error)
	GetListener(ListenerConfig) (EventListener, error)
}

// EmiterConfig holds data for configuring emiters
type EmiterConfig struct {
	// Only for AMQP, name of the exchange to which messages will be published
	ExchangeName string

	// allows to log to specific source
	Logger *log.Logger
}

// ListenerConfig holds data for configuring listeners
type ListenerConfig struct {
	// used for naming kafka consumerGroup or AMQP's queue
	ClientName string

	// if Broadcast is true, individual listener is returned
	Broadcast bool

	// List of events listener should interject
	Events []Event

	// default: msgqueue.jsonDecoder
	Decoder *Decoder

	// Only for Kafka, default: OffsetNewest (-1), 0 leaves default option, to start from first message written use OffsetOldest or -2
	Offset int64

	// Only for Kafka, allows to set number of partitions for a topic, maps 'topic' -> no. of partitions
	SetPartitions map[string]int32

	// allows to log to specific source
	Logger *log.Logger
}
