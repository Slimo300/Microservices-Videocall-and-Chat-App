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
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/cache"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/routes"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/ws"
)

func main() {

	config, err := configuration.LoadConfig(os.Getenv("CHAT_CONFIG"))
	if err != nil {
		log.Fatalf("Error when loading configuration: %v", err)
	}

	db, err := database.Setup(config.WSService.DBType, config.WSService.DBAddress)
	if err != nil {
		log.Fatalf("Error when connecting to database: %v", err)
	}
	tokenService, err := auth.NewGRPCTokenClient(config.AuthAddress)
	if err != nil {
		log.Fatalf("Error when connecting to token service: %v", err)
	}

	conf := sarama.NewConfig()
	conf.ClientID = "websocketService"
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
		reflect.TypeOf(events.GroupDeletedEvent{}),
		reflect.TypeOf(events.MemberCreatedEvent{}),
		reflect.TypeOf(events.MemberDeletedEvent{}),
		reflect.TypeOf(events.MemberUpdatedEvent{}),
		reflect.TypeOf(events.MessageDeletedEvent{}),
		reflect.TypeOf(events.InviteSentEvent{}),
		reflect.TypeOf(events.InviteRespondedEvent{}),
	); err != nil {
		log.Fatal(err)
	}
	listener, err := kafka.NewKafkaEventListener(client, mapper, kafka.KafkaTopic{Name: "messages"}, kafka.KafkaTopic{Name: "groups"})
	if err != nil {
		log.Fatal(err)
	}

	messageChan := make(chan *ws.Message)
	actionChan := make(chan msgqueue.Event)
	server := &handlers.Server{
		DB:           db,
		CodeCache:    cache.NewCache(5 * time.Second),
		TokenService: tokenService,
		Emitter:      emitter,
		Listener:     listener,
		Hub:          ws.NewHub(messageChan, actionChan, config.Origin),
		MessageChan:  messageChan,
		EventChan:    actionChan,
	}
	go server.RunHub()
	handler := routes.Setup(server, config.Origin)

	go server.RunListener()

	httpServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", config.WSService.HTTPPort),
	}
	httpsServer := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", config.WSService.HTTPSPort),
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
