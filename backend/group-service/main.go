package main

import (
	"context"
	"errors"
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
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	db, err := orm.Setup()
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.Setup()
	authService, err := auth.NewGRPCTokenClient(":9000")
	if err != nil {
		panic("Couldn't connect to grpc auth server")
	}
	server := handlers.NewServer(db, &storage, authService)
	routes.Setup(engine, server)
	go server.RunHub()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

}
