package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/configuration"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/routes"
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
	tokenService, err := auth.NewGRPCTokenClient(config.TokenService.GRPCPort)
	if err != nil {
		log.Fatalf("Error when connecting to token service: %v", err)
	}
	server := &handlers.Server{DB: db, TokenService: tokenService}
	routes.Setup(engine, server)

	httpServer := &http.Server{
		Handler: engine,
		Addr:    config.WSService.HTTPPort,
	}
	httpsServer := &http.Server{
		Handler: engine,
		Addr:    config.WSService.HTTPSPort,
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
