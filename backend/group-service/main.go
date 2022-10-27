package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database/orm"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/routes"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/configuration"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	config, err := configuration.LoadConfig(os.Getenv("CHAT_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := orm.Setup(config.GroupService.DBAddress)
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.Setup()
	authService, err := auth.NewGRPCTokenClient(config.TokenService.GRPCPort)
	if err != nil {
		panic("Couldn't connect to grpc auth server")
	}
	server := handlers.NewServer(db, &storage, authService)
	routes.Setup(engine, server)

	httpServer := &http.Server{
		Handler: engine,
		Addr:    config.GroupService.HTTPPort,
	}
	httpsServer := &http.Server{
		Handler: engine,
		Addr:    config.GroupService.HTTPSPort,
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
