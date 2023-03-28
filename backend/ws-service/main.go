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

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	"github.com/Slimo300/chat-wsservice/internal/cache"
	"github.com/Slimo300/chat-wsservice/internal/config"
	"github.com/Slimo300/chat-wsservice/internal/database"
	"github.com/Slimo300/chat-wsservice/internal/eventprocessor"
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

	emiter, listener, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatalf("Error setting up kafka: %v", err)
	}

	messageChan := make(chan *ws.Message)
	actionChan := make(chan msgqueue.Event)

	eventProcessor := eventprocessor.NewEventProcessor(db, listener, actionChan)
	go eventProcessor.ProcessEvents()

	server := &handlers.Server{
		DB:          db,
		CodeCache:   cache.NewCache(5 * time.Second),
		TokenClient: tokenClient,
		Emitter:     emiter,
		Listener:    listener,
		Hub:         ws.NewHub(messageChan, actionChan, conf.Origin),
		MessageChan: messageChan,
	}

	go server.RunHub()
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
