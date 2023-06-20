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

	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/config"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/database/elastic"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/eventprocessor"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/routes"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
)

func main() {
	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Couln't load config: %v", err)
	}

	tokenClient, err := auth.NewGRPCTokenClient(conf.TokenServiceAddress)
	if err != nil {
		log.Fatalf("Error connecting to token service: %v", err)
	}

	listener, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatalf("Error creating kafka listener: %v", err)
	}

	es, err := elastic.NewElasticSearchDB([]string{conf.DBAddress}, conf.DBUser, conf.DBPassword)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializng new event processor")

	eventProcessor := eventprocessor.NewEventProcessor(es, listener)
	go eventProcessor.ProcessEvents("users")

	server := handlers.Server{
		DB:          es,
		Listener:    listener,
		TokenClient: tokenClient,
	}

	handler := routes.Setup(&server, conf.Origin)

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
