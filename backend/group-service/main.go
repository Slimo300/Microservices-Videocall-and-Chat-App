package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database/orm"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/routes"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/configuration"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/kafka"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
)

func main() {

	config, err := configuration.LoadConfig(os.Getenv("CHAT_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := orm.Setup(config.GroupService.DBType, config.GroupService.DBAddress)
	if err != nil {
		log.Fatal(err)
	}
	storage, err := storage.NewS3Storage(config.S3Bucket, config.Origin)
	if err != nil {
		log.Fatalf("Error connecting to AWS S3: %v", err)
	}
	authService, err := auth.NewGRPCTokenClient(config.AuthAddress)
	if err != nil {
		log.Fatalf("Couldn't connect to grpc auth server: %v", err)
	}

	conf := sarama.NewConfig()
	conf.ClientID = "groupsService"
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
	mapper := msgqueue.NewDynamicEventMapper()
	if err := mapper.RegisterTypes(
		reflect.TypeOf(events.UserRegisteredEvent{}),
		reflect.TypeOf(events.UserPictureModifiedEvent{}),
	); err != nil {
		log.Fatal(err)
	}
	listener, err := kafka.NewKafkaEventListener(client, mapper, kafka.KafkaTopic{Name: "users"})
	if err != nil {
		log.Fatal(err)
	}

	server := handlers.Server{
		DB:           db,
		Storage:      storage,
		TokenService: authService,
		Emitter:      emitter,
		Listener:     listener,
		MaxBodyBytes: 4194304,
	}
	handler := routes.Setup(&server, config.Origin)

	go server.RunListener()

	httpServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", config.GroupService.HTTPPort),
	}
	httpsServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", config.GroupService.HTTPSPort),
	}

	errChan := make(chan error)

	go func() {
		errChan <- httpsServer.ListenAndServeTLS(config.Certificate, config.PrivKeyFile)
	}()
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
