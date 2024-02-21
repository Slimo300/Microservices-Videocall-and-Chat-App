package main

import (
	"context"
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
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-gateway-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-gateway-service/database/redis"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-gateway-service/eventprocessor"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-gateway-service/handlers"
)

func main() {
	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Error reading config from environment variables: %v", err)
	}

	db, err := redis.Setup(conf.DBAddress, conf.DBPassword)
	if err != nil {
		log.Fatalf("Error when trying to connect to Redis: %v", err)
	}

	builder, err := builder.NewBrokerBuilder(msgqueue.ParseBrokerType(conf.BrokerType), conf.BrokerAddress)
	if err != nil {
		log.Fatalf("Error creating broker builder: %v", err)
	}

	listener, err := builder.GetListener(msgqueue.ListenerConfig{
		ClientName: "webrtc-gateway-service",
		Events: []msgqueue.Event{
			events.ServiceStartedEvent{},
		},
	})
	if err != nil {
		log.Fatalf("Error building listener: %v", err)
	}

	go eventprocessor.NewEventProcessor(listener, db).ProcessEvents("webrtc")

	handler := handlers.Setup(db, conf.Origin)

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
