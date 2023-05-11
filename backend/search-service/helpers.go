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

// startHTTPServer runs HTTPS Server if SSL certificate is provided
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

// kafkaSetup starts kafka listener
func kafkaSetup(brokerAddresses []string) (msgqueue.EventListener, error) {

	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "searchService"
	brokerConf.Version = sarama.V2_3_0_0
	client, err := sarama.NewClient(brokerAddresses, brokerConf)
	if err != nil {
		return nil, err
	}

	mapper := msgqueue.NewDynamicEventMapper()
	if err := mapper.RegisterTypes(
		reflect.TypeOf(events.UserRegisteredEvent{}),
		reflect.TypeOf(events.UserPictureModifiedEvent{}),
	); err != nil {
		return nil, err
	}

	listener, err := kafka.NewConsumerGroupEventListener(client, "search-service", mapper, &kafka.ListenerOptions{
		Logger: log.New(os.Stdout, "[listener]: ", log.Flags()),
	})
	if err != nil {
		return nil, err
	}

	return listener, nil
}
