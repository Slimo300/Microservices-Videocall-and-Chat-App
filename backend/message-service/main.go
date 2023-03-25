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
	"github.com/Slimo300/chat-messageservice/internal/config"
	"github.com/Slimo300/chat-messageservice/internal/database/orm"
	"github.com/Slimo300/chat-messageservice/internal/handlers"
	"github.com/Slimo300/chat-messageservice/internal/routes"
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/kafka"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
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

	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "messagesService"
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
		reflect.TypeOf(events.MessageSentEvent{}),
	); err != nil {
		log.Fatal(err)
	}
	listener, err := kafka.NewKafkaEventListener(client, mapper, kafka.KafkaTopic{Name: "wsmessages"}, kafka.KafkaTopic{Name: "groups"})
	if err != nil {
		log.Fatal(err)
	}

	tokenClient, err := tokens.NewGRPCTokenClient(conf.TokenServiceAddress)
	if err != nil {
		log.Fatalf("Couldn't connect to grpc auth server: %v", err)
	}
	storage, err := storage.NewS3Storage(conf.S3Bucket, conf.Origin)
	if err != nil {
		log.Fatalf("Couldn't establish s3 session: %v", err)
	}

	server := &handlers.Server{
		DB:          db,
		TokenClient: tokenClient,
		Emitter:     emitter,
		Listener:    listener,
		Storage:     storage,
	}
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
