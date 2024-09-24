package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/builder"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database/orm"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/eventprocessor"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/routes"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/storage/s3"
)

func getPublicKey() (*rsa.PublicKey, error) {

	bytePubKey, err := os.ReadFile("/rsa/public.key")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytePubKey)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func main() {

	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	pubKey, err := getPublicKey()
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	db, err := orm.Setup(conf.DBAddress)
	if err != nil {
		log.Fatal(err)
	}

	builder, err := builder.NewBrokerBuilder(msgqueue.ParseBrokerType(conf.BrokerType), conf.BrokerAddress)
	if err != nil {
		log.Fatalf("Error when creating broker builder: %v", err)
	}

	emiter, err := builder.GetEmiter(msgqueue.EmiterConfig{
		ExchangeName: "message",
	})
	if err != nil {
		log.Fatalf("Error when building emitter: %v", err)
	}

	listener, err := builder.GetListener(msgqueue.ListenerConfig{
		ClientName: "message-service",
		Events: []msgqueue.Event{
			events.GroupDeletedEvent{},
			events.MemberCreatedEvent{},
			events.MemberDeletedEvent{},
			events.MemberUpdatedEvent{},
			events.MessageSentEvent{},
		},
	})
	if err != nil {
		log.Fatalf("Error when building listener: %v", err)
	}

	storage, err := s3.NewS3Storage(conf.StorageKeyID, conf.StorageKeySecret, conf.Bucket, s3.WithCORS(conf.Origin))
	if err != nil {
		log.Fatalf("Couldn't establish s3 session: %v", err)
	}

	go eventprocessor.NewEventProcessor(listener, db, storage).ProcessEvents("wsmessage", "group")

	server := handlers.NewServer(db, pubKey, emiter, storage)
	handler := routes.Setup(server, conf.Origin)

	httpServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", conf.HTTPPort),
	}

	errChan := make(chan error)

	go func() { errChan <- httpServer.ListenAndServe() }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v\n", err)
		}
	case err := <-errChan:
		log.Fatal(err)
	}
}
