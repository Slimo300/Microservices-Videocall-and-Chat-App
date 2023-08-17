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

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/config"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/database/redis"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/eventprocessor"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/webrtc"
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
		log.Fatalf("Error reading configuration: %v", err)
	}

	pubKey, err := getPublicKey()
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	db, err := redis.Setup(conf.DBAddress, conf.DBPassword)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	_, dbListener, relayListener, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatalf("Error setting up kafka: %v", err)
	}

	relayChan := make(chan msgqueue.Event)

	go eventprocessor.NewDBEventProcessor(dbListener, db).ProcessEvents("groups")
	go eventprocessor.NewRelayEventProcessor(relayListener, relayChan).ProcessEvents("groups", "webrtc")

	server := &handlers.Server{
		DB:        db,
		PublicKey: pubKey,
		Relay:     webrtc.NewRoomsRelay(),
	}

	handler := server.Setup(conf.Origin)

	httpServer := http.Server{
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
