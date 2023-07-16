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

	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/config"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database/orm"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/eventprocessor"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/routes"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/storage"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
)

func main() {

	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	db, err := orm.Setup(conf.DBAddress)
	if err != nil {
		log.Fatal(err)
	}

	emiter, listener, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatal(err)
	}

	tokenClient, err := auth.NewGRPCTokenClient(conf.TokenServiceAddress)
	if err != nil {
		log.Fatalf("Couldn't connect to grpc auth server: %v", err)
	}

	// storage, err := storage.NewS3Storage(conf.S3Bucket, conf.Origin)
	// if err != nil {
	// 	log.Fatalf("Couldn't establish s3 session: %v", err)
	// }

	go eventprocessor.NewEventProcessor(listener, db, new(storage.MockStorage)).ProcessEvents("wsmessages", "groups")

	server := handlers.NewServer(db, tokenClient, emiter, new(storage.MockStorage))
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
