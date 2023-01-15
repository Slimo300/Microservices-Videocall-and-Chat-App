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
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/configuration"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/kafka"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/database/elastic"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/routes"
)

func main() {
	conf, err := configuration.LoadConfig(os.Getenv("CHAT_CONFIG"))
	if err != nil {
		log.Fatalf("Couln't load configuration file: %v", err)
	}

	authService, err := auth.NewGRPCTokenClient(conf.AuthAddress)
	if err != nil {
		log.Fatalf("Error connecting to token service: %v", err)
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = "searchService"
	saramaConfig.Version = sarama.V2_3_0_0
	client, err := sarama.NewClient(conf.BrokersAddresses, saramaConfig)
	if err != nil {
		log.Fatalf("Error when connecting to kafka: %v", err)
	}

	mapper := msgqueue.NewDynamicEventMapper()
	if err := mapper.RegisterEventType(reflect.TypeOf(events.UserRegisteredEvent{})); err != nil {
		log.Fatalf("Error when registering type: %v", err)
	}

	listener, err := kafka.NewKafkaEventListener(client, mapper, kafka.KafkaTopic{Name: "users"})
	if err != nil {
		log.Fatalf("Error creating kafka listener: %v", err)
	}
	listener.Offset = sarama.OffsetOldest

	es, err := elastic.NewElasticSearchDB(conf.SearchService.Addresses, conf.SearchService.Username, conf.SearchService.Password)
	if err != nil {
		log.Fatal(err)
	}

	server := handlers.Server{
		DB:          es,
		Listener:    listener,
		AuthService: authService,
	}
	go server.RunListener()

	handler := routes.Setup(&server, conf.Origin)

	httpServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", conf.SearchService.HTTPPort),
	}
	httpsServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", conf.SearchService.HTTPSPort),
	}

	errChan := make(chan error)

	go func() {
		errChan <- httpsServer.ListenAndServeTLS(conf.Certificate, conf.PrivKeyFile)
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
