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

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/email"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database/orm"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/routes"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/storage"
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
	tokenConn, err := grpc.Dial(conf.TokenServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to token service")
	}
	tokenClient := auth.NewTokenServiceClient(tokenConn)

	emailConn, err := grpc.Dial(conf.EmailServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to token service")
	}
	emailClient := email.NewEmailServiceClient(emailConn)

	emiter, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatalf("Error setting up kafka: %v", err)
	}

	// Setup for handling image uploads to s3 and email sending
	storage, err := storage.NewS3Storage(conf.StorageKeyID, conf.StorageKeySecret, conf.Bucket)
	if err != nil {
		log.Fatal(err)
	}

	server := &handlers.Server{
		DB:           db,
		TokenClient:  tokenClient,
		EmailClient:  emailClient,
		Emitter:      emiter,
		TokenKey:     pubkey,
		ImageStorage: storage,
		MaxBodyBytes: 4194304, //4MB
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
