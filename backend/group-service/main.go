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

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/app"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database/orm"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/eventprocessor"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/amqp"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage/s3"
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
		log.Fatalf("Couldn't read config: %v", err)
	}

	pubKey, err := getPublicKey()
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	db, err := orm.NewGroupsGormRepository(conf.DBAddress)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := s3.NewS3Storage(context.Background(), conf.StorageKeyID, conf.StorageKeySecret, conf.Bucket, s3.WithRegion(conf.StorageRegion))
	if err != nil {
		log.Fatalf("Error connecting to AWS S3: %v", err)
	}

	amqpBuilder, err := amqp.NewAMQPBuilder(conf.BrokerAddress)
	if err != nil {
		log.Fatalf("Error creating broker builder: %v", err)
	}

	emiter, err := amqpBuilder.GetEmiter(msgqueue.EmiterConfig{
		ExchangeName: "group",
	})
	if err != nil {
		log.Fatalf("Error when building emitter: %v", err)
	}

	listener, err := amqpBuilder.GetListener(msgqueue.ListenerConfig{
		ClientName: "group-service",

		Events: []msgqueue.Event{
			events.UserVerifiedEvent{},
			events.UserPictureModifiedEvent{},
		},
	})
	if err != nil {
		log.Fatalf("Error when building listener: %v", err)
	}

	app := app.NewApplication(db, storage, emiter)

	go eventprocessor.NewEventProcessor(app, listener).ProcessEvents("user")

	httpServer := &http.Server{
		Handler: handlers.NewServer(app, pubKey, conf.Origin),
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
