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

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	db, err := database.Setup()
	if err != nil {
		log.Fatal(err)
	}
	tokenService, err := auth.NewGRPCTokenClient(":9000")
	if err != nil {
		panic("Couldn't connect to grpc auth server")
	}
	server := &handlers.Server{DB: db, TokenService: tokenService}
	routes.Setup(engine, server)

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
