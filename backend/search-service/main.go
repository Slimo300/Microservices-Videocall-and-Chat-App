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

	"github.com/Slimo300/chat-searchservice/internal/config"
	"github.com/Slimo300/chat-searchservice/internal/database/elastic"
	"github.com/Slimo300/chat-searchservice/internal/eventprocessor"
	"github.com/Slimo300/chat-searchservice/internal/handlers"
	"github.com/Slimo300/chat-searchservice/internal/routes"

	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
)

func main() {
	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Couln't load config: %v", err)
	}

	tokenClient, err := tokens.NewGRPCTokenClient(conf.TokenServiceAddress)
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

	eventProcessor := eventprocessor.NewEventProcessor(listener, es)
	go eventProcessor.ProcessEvents()

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
	httpsServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", conf.HTTPSPort),
	}

	errChan := make(chan error)

	go startHTTPSServer(httpsServer, conf.CertDir, errChan)
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
		if err := httpsServer.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v\n", err)
		}
	case err := <-errChan:
		log.Fatal(err)
	}
}
