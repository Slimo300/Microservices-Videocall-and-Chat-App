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
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/config"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/database/redis"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/eventprocessor"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/routes"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/ws"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
)

func main() {

	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Error when loading configuration: %v", err)
	}

	db, err := redis.Setup(conf.DBAddress, conf.DBPassword)
	if err != nil {
		log.Fatalf("Error when connecting to database: %v", err)
	}
	tokenClient, err := auth.NewGRPCTokenClient(conf.TokenServiceAddress)
	if err != nil {
		log.Fatalf("Error when connecting to token service: %v", err)
	}

	emiter, dbListener, hubListener, err := kafkaSetup([]string{conf.BrokerAddress})
	if err != nil {
		log.Fatalf("Error setting up kafka: %v", err)
	}

	messageChan := make(chan *ws.Message)
	actionChan := make(chan msgqueue.Event)

	go eventprocessor.NewDBEventProcessor(dbListener, db).ProcessEvents("groups")
	go eventprocessor.NewHubEventProcessor(hubListener, actionChan).ProcessEvents("groups", "messages", "wsmessages")

	server := &handlers.Server{
		DB:          db,
		TokenClient: tokenClient,
		Emitter:     emiter,
		Hub:         ws.NewHub(messageChan, actionChan, conf.Origin),
		MessageChan: messageChan,
	}

	go server.RunHub()
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
