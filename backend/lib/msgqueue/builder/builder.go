package builder

import (
	"fmt"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/amqp"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/kafka"
)

// NewBrokerBuilder creates a builder for emiters and listeners either for kafka or for amqp
func NewBrokerBuilder(brokerType msgqueue.BrokerType, brokerAddress string) (msgqueue.BrokerBuilder, error) {

	switch brokerType {
	case msgqueue.AMQP_BROKER:
		return amqp.NewAMQPBuilder(brokerAddress)
	case msgqueue.KAFKA_BROKER:
		return kafka.NewKafkaBuilder([]string{brokerAddress})
	default:
		return nil, fmt.Errorf("error creating builder: invalid broker type")
	}
}
