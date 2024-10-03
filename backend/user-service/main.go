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

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/handlers"
)

func readPublicKey() (*rsa.PublicKey, error) {

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
	pubkey, err := readPublicKey()
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Error when loading configuration: %v", err)
	}

	app := app.NewApplication(conf)
	handler := handlers.NewServer(app, pubkey, conf.Domain, conf.Origin, 4000000)

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
