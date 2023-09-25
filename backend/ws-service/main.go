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

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
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

	emiter, dbListener, hubListener, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatalf("Error setting up kafka: %v", err)
	}

	eventChan := make(chan msgqueue.Event)

	go eventprocessor.NewDBEventProcessor(dbListener, db).ProcessEvents("groups")
	go eventprocessor.NewHubEventProcessor(hubListener, eventChan).ProcessEvents("groups", "messages", "wsmessages")

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
