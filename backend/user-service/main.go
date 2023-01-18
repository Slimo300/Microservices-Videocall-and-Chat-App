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

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/configuration"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/kafka"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/database/orm"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/email"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/routes"
)

func main() {

	config, err := configuration.LoadConfig(os.Getenv("CHAT_CONFIG"))
	if err != nil {
		log.Fatalf("Error when loading configuration: %v", err)
	}
	// Setting up MySQL connection
	db, err := orm.Setup(config.UserService.DBType, config.UserService.DBAddress, orm.WithConfig(orm.DBConfig{
		VerificationCodeDuration: 24 * time.Hour,
		ResetCodeDuration:        10 * time.Minute,
	}))
	if err != nil {
		log.Fatalf("Error when connecting to database: %v", err)
	}

	// connecting to authentication server
	tokenService, err := auth.NewGRPCTokenClient(config.AuthAddress)
	if err != nil {
		log.Fatalf("Error when connecting to token service: %v", err)
	}

	// kafka broker setup
	conf := sarama.NewConfig()
	conf.ClientID = "userService"
	conf.Version = sarama.V2_3_0_0
	conf.Producer.Return.Successes = true
	client, err := sarama.NewClient(config.BrokersAddresses, conf)
	if err != nil {
		log.Fatal(err)
	}

	emitter, err := kafka.NewKafkaEventEmiter(client)
	if err != nil {
		log.Fatal(err)
	}

	// Setup for handling image uploads to s3 and email sending
	storage, err := storage.Setup(config.S3Bucket)
	if err != nil {
		log.Fatal(err)
	}
	emailService, err := email.NewSMTPService(config.EmailTemplateDir, config.EmailFrom, config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass)
	if err != nil {
		log.Fatal(err)
	}

	server := &handlers.Server{
		Origin:       config.Origin,
		DB:           db,
		TokenService: tokenService,
		Emitter:      emitter,
		ImageStorage: storage,
		EmailService: emailService,
		MaxBodyBytes: 4194304,
		Domain:       config.Domain,
	}
	handler := routes.Setup(server, config.Origin)

	httpServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", config.UserService.HTTPPort),
	}
	httpsServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", config.UserService.HTTPSPort),
	}

	errChan := make(chan error)

	go func() {
		errChan <- httpsServer.ListenAndServeTLS(config.Certificate, config.PrivKeyFile)
	}()
	go func() { errChan <- httpServer.ListenAndServe() }()

	quit := make(chan os.Signal)
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
