package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/kafka"
	"github.com/Slimo300/chat-wsservice/internal/cache"
	"github.com/Slimo300/chat-wsservice/internal/config"
	"github.com/Slimo300/chat-wsservice/internal/database"
	"github.com/Slimo300/chat-wsservice/internal/handlers"
	"github.com/Slimo300/chat-wsservice/internal/routes"
	"github.com/Slimo300/chat-wsservice/internal/ws"

	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
)

func main() {

	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Error when loading configuration: %v", err)
	}

	db, err := database.Setup(conf.DBAddress)
	if err != nil {
		log.Fatalf("Error when connecting to database: %v", err)
	}
	tokenClient, err := tokens.NewGRPCTokenClient(conf.TokenServiceAddress)
	if err != nil {
		log.Fatalf("Error when connecting to token service: %v", err)
	}

	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "websocketService"
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
		DB:          db,
		CodeCache:   cache.NewCache(5 * time.Second),
		TokenClient: tokenClient,
		Emitter:     emitter,
		Listener:    listener,
		Hub:         ws.NewHub(messageChan, actionChan, conf.Origin),
		MessageChan: messageChan,
		EventChan:   actionChan,
	}
	go server.RunHub()
	handler := routes.Setup(server, conf.Origin)

	go server.RunListener()

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
