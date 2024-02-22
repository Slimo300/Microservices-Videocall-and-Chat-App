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
	"strings"
	"syscall"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/builder"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/database/redis"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/eventprocessor"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/routes"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/ws-service/ws"
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
		log.Fatalf("Error when loading configuration: %v", err)
	}

	db, err := redis.Setup(conf.DBAddress, conf.DBPassword)
	if err != nil {
		log.Fatalf("Error when connecting to database: %v", err)
	}

	pubKey, err := getPublicKey()
	if err != nil {
		log.Fatalf("Error when reading public key: %v", err)
	}

	builder, err := builder.NewBrokerBuilder(msgqueue.ParseBrokerType(conf.BrokerType), conf.BrokerAddress)
	if err != nil {
		log.Fatalf("Error when creating broker builder: %v", err)
	}

	emiter, err := builder.GetEmiter(msgqueue.EmiterConfig{
		ExchangeName: "wsmessage",
	})
	if err != nil {
		log.Fatalf("Error when building emitter: %v", err)
	}

	dbListener, err := builder.GetListener(msgqueue.ListenerConfig{
		ClientName: "ws-service",
		Events: []msgqueue.Event{
			events.GroupDeletedEvent{},
			events.MemberCreatedEvent{},
			events.MemberDeletedEvent{},
			events.MemberUpdatedEvent{},
		},
	})
	if err != nil {
		log.Fatalf("Error when building data replication listener: %v", err)
	}

	podName := strings.Split(conf.PodName, "-")
	if len(podName) < 4 {
		log.Fatalf("Invalid Pod name %v", podName)
	}

	hubListener, err := builder.GetListener(msgqueue.ListenerConfig{
		ClientName: "ws-service-" + podName[2] + podName[3],
		Broadcast:  true,
		Events: []msgqueue.Event{
			events.GroupDeletedEvent{},
			events.InviteSentEvent{},
			events.InviteRespondedEvent{},
			events.MemberCreatedEvent{},
			events.MemberDeletedEvent{},
			events.MemberUpdatedEvent{},
			events.MessageDeletedEvent{},
			events.MessageSentEvent{},
		},
	})
	if err != nil {
		log.Fatalf("Error when building hub listener: %v", err)
	}

	eventChan := make(chan msgqueue.Event)

	go eventprocessor.NewDBEventProcessor(dbListener, db).ProcessEvents("group")
	go eventprocessor.NewHubEventProcessor(hubListener, eventChan).ProcessEvents("group", "message", "wsmessage")

	server := &handlers.Server{
		DB:        db,
		PublicKey: pubKey,
		Hub:       ws.NewHub(eventChan, emiter, conf.Origin),
	}

	go server.Hub.Run()
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
