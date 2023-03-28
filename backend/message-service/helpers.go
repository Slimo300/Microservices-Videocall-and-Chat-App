package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/kafka"
)

// starts HTTPS server if SSL certificate is provided
func startHTTPSServer(httpsServer *http.Server, certDir string, errChan chan<- error) {
	cert := filepath.Join(certDir, "cert.pem")
	if _, err := os.Stat(cert); err != nil {
		log.Printf("Couldn't start https server. No cert.pem or key.pem in %s\n", certDir)
		return
	}
	key := filepath.Join(certDir, "key.pem")
	if _, err := os.Stat(key); err != nil {
		log.Printf("Couldn't start https server. No cert.pem or key.pem in %s\n", certDir)
		return
	}

	log.Printf("HTTPS Server starting on: %s", httpsServer.Addr)
	errChan <- httpsServer.ListenAndServeTLS(cert, key)
}

// returns kafka EventEmitter and EventListener
func kafkaSetup(brokerAddresses []string) (msgqueue.EventEmiter, msgqueue.EventListener, error) {
	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "messagesService"
	brokerConf.Version = sarama.V2_3_0_0
	brokerConf.Producer.Return.Successes = true
	client, err := sarama.NewClient(brokerAddresses, brokerConf)
	if err != nil {
		return nil, nil, err
	}

	emiter, err := kafka.NewKafkaEventEmiter(client)
	if err != nil {
		return nil, nil, err
	}
	mapper := msgqueue.NewDynamicEventMapper()
	if err := mapper.RegisterTypes(
		reflect.TypeOf(events.GroupDeletedEvent{}),
		reflect.TypeOf(events.MemberCreatedEvent{}),
		reflect.TypeOf(events.MemberDeletedEvent{}),
		reflect.TypeOf(events.MemberUpdatedEvent{}),
		reflect.TypeOf(events.MessageSentEvent{}),
	); err != nil {
		return nil, nil, err
	}
	listener, err := kafka.NewKafkaEventListener(client, mapper, kafka.KafkaTopic{Name: "wsmessages"}, kafka.KafkaTopic{Name: "groups"})
	if err != nil {
		return nil, nil, err
	}

	return emiter, listener, nil
}
