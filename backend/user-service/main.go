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

	emails "github.com/Slimo300/chat-emailservice/pkg/client"
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"

	"github.com/Slimo300/chat-userservice/internal/config"
	"github.com/Slimo300/chat-userservice/internal/database/orm"
	"github.com/Slimo300/chat-userservice/internal/handlers"
	"github.com/Slimo300/chat-userservice/internal/routes"
	"github.com/Slimo300/chat-userservice/internal/storage"
)

func main() {

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
	tokenClient, err := tokens.NewGRPCTokenClient(conf.TokenServiceAddress)
	if err != nil {
		log.Fatalf("Error when connecting to token service: %v", err)
	}

	emiter, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatalf("Error setting up kafka: %v", err)
	}

	// Setup for handling image uploads to s3 and email sending
	storage, err := storage.NewS3Storage(conf.S3Bucket, conf.Origin)
	if err != nil {
		log.Fatal(err)
	}
	emailClient, err := emails.NewGRPCEmailClient(conf.EmailServiceAddress)
	if err != nil {
		log.Fatal(err)
	}

	server := &handlers.Server{
		DB:           db,
		TokenClient:  tokenClient,
		EmailClient:  emailClient,
		Emitter:      emiter,
		ImageStorage: storage,
		MaxBodyBytes: 4194304,
		Domain:       conf.Domain,
	}
	handler := routes.Setup(server, conf.Origin)

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
