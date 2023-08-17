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

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/email"

	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/config"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/database/orm"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/routes"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/storage"
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
	// Setting up MySQL connection
	db, err := orm.Setup(conf.DBAddress, orm.WithConfig(orm.DBConfig{
		VerificationCodeDuration: 24 * time.Hour,
		ResetCodeDuration:        10 * time.Minute,
	}))
	if err != nil {
		log.Fatalf("Error when connecting to database: %v", err)
	}

	// connecting to authentication server
	tokenClient, err := auth.NewGRPCTokenClient(conf.TokenServiceAddress)
	if err != nil {
		log.Fatalf("Error when connecting to token service: %v", err)
	}

	emiter, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatalf("Error setting up kafka: %v", err)
	}

	// Setup for handling image uploads to s3 and email sending
	// storage, err := storage.NewS3Storage(conf.S3Bucket, conf.Origin)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	emailClient, err := email.NewGRPCEmailClient(conf.EmailServiceAddress)
	if err != nil {
		log.Fatal(err)
	}

	server := &handlers.Server{
		DB:           db,
		TokenClient:  tokenClient,
		EmailClient:  emailClient,
		Emitter:      emiter,
		TokenKey:     pubkey,
		ImageStorage: new(storage.MockStorage),
		MaxBodyBytes: 4194304,
		Domain:       conf.Domain,
	}
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
