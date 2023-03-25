package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/kafka"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	emails "github.com/Slimo300/chat-emailservice/pkg/client"
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
	"github.com/Slimo300/chat-userservice/internal/config"
	"github.com/Slimo300/chat-userservice/internal/database/orm"
	"github.com/Slimo300/chat-userservice/internal/handlers"
	"github.com/Slimo300/chat-userservice/internal/routes"
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

	// kafka broker setup
	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "userService"
	brokerConf.Version = sarama.V2_3_0_0
	brokerConf.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{conf.BrokerAddress}, brokerConf)
	if err != nil {
		log.Fatal(err)
	}

	emitter, err := kafka.NewKafkaEventEmiter(client)
	if err != nil {
		log.Fatal(err)
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
		Emitter:      emitter,
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

func startHTTPSServer(httpsServer *http.Server, certDir string, errChan chan<- error) {
	cert := filepath.Join(certDir, "cert.pem")
	log.Println(cert)
	if _, err := os.Stat(cert); err != nil {
		log.Printf("Couldn't start https server. No cert.pem or key.pem in %s\n", certDir)
		return
	}

	key := filepath.Join(certDir, "key.pem")
	if _, err := os.Stat(key); err != nil {
		log.Printf("Couldn't start https server. No cert.pem or key.pem in %s\n", certDir)
		return
	}

	log.Printf("HTTPS Server starting on: %s", httpsServer.Addr)
	errChan <- httpsServer.ListenAndServeTLS(cert, key)
}
