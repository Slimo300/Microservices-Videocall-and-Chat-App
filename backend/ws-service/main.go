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
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/routes"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	config, err := configuration.LoadConfig(os.Getenv("CHAT_CONFIG"))
	if err != nil {
		log.Fatalf("Error when loading configuration: %v", err)
	}

	db, err := database.Setup(config.WSService.DBAddress)
	if err != nil {
		log.Fatalf("Error when connecting to database: %v", err)
	}
	tokenService, err := auth.NewGRPCTokenClient(fmt.Sprintf(":%s", config.TokenService.GRPCPort))
	if err != nil {
		log.Fatalf("Error when connecting to token service: %v", err)
	}

	conf := sarama.NewConfig()
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
		reflect.TypeOf(events.GroupDeletedEvent{}),
		reflect.TypeOf(events.MemberCreatedEvent{}),
		reflect.TypeOf(events.MemberDeletedEvent{}),
		reflect.TypeOf(events.MemberUpdatedEvent{}),
		reflect.TypeOf(events.MessageSentEvent{}),
	); err != nil {
		log.Fatal(err)
	}
	listener, err := kafka.NewKafkaEventListener(client, mapper, kafka.KafkaTopic{Name: "messages"}, kafka.KafkaTopic{Name: "groups"})
	if err != nil {
		log.Fatal(err)
	}

	messageChan := make(chan *ws.Message)
	server := &handlers.Server{
		DB:           db,
		TokenService: tokenService,
		Emitter:      emitter,
		Listener:     listener,
		Hub:          ws.NewHub(messageChan),
		MessageChan:  messageChan,
	}
	routes.Setup(engine, server)

	go server.RunListener()

	httpServer := &http.Server{
		Handler: engine,
		Addr:    fmt.Sprintf(":%s", config.WSService.HTTPPort),
	}
	httpsServer := &http.Server{
		Handler: engine,
		Addr:    fmt.Sprintf(":%s", config.WSService.HTTPSPort),
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
