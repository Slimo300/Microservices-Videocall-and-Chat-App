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
	"github.com/Slimo300/chat-searchservice/internal/config"
	"github.com/Slimo300/chat-searchservice/internal/database/elastic"
	"github.com/Slimo300/chat-searchservice/internal/handlers"
	"github.com/Slimo300/chat-searchservice/internal/routes"
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
)

func main() {
	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Couln't load config: %v", err)
	}

	tokenClient, err := tokens.NewGRPCTokenClient(conf.TokenServiceAddress)
	if err != nil {
		log.Fatalf("Error connecting to token service: %v", err)
	}

	brokerConf := sarama.NewConfig()
	brokerConf.ClientID = "searchService"
	brokerConf.Version = sarama.V2_3_0_0
	client, err := sarama.NewClient([]string{conf.BrokerAddress}, brokerConf)
	if err != nil {
		log.Fatalf("Error when connecting to kafka: %v", err)
	}

	mapper := msgqueue.NewDynamicEventMapper()
	if err := mapper.RegisterTypes(
		reflect.TypeOf(events.UserRegisteredEvent{}),
		reflect.TypeOf(events.UserPictureModifiedEvent{}),
	); err != nil {
		log.Fatalf("Error when registering types by mapper: %v", err)
	}

	listener, err := kafka.NewKafkaEventListener(client, mapper, kafka.KafkaTopic{Name: "users"})
	if err != nil {
		log.Fatalf("Error creating kafka listener: %v", err)
	}

	es, err := elastic.NewElasticSearchDB([]string{conf.DBAddress}, conf.DBUser, conf.DBPassword)
	if err != nil {
		log.Fatal(err)
	}

	server := handlers.Server{
		DB:          es,
		Listener:    listener,
		TokenClient: tokenClient,
	}
	go server.RunListener()

	handler := routes.Setup(&server, conf.Origin)

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
