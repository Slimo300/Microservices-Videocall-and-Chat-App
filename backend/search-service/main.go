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

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/database/elastic"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/eventprocessor"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/search-service/routes"
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
		log.Fatalf("Couln't load config: %v", err)
	}

	pubKey, err := getPublicKey()
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
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
		PublicKey: pubKey,
		DB:        es,
		Listener:  listener,
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
