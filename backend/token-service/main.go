package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Slimo300/chat-tokenservice/internal/config"
	"github.com/Slimo300/chat-tokenservice/internal/handlers"
	"github.com/Slimo300/chat-tokenservice/internal/repo/redis"
	"github.com/Slimo300/chat-tokenservice/pkg/client/pb"
	"google.golang.org/grpc"
)

func main() {

	config, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Couldn't read configuration: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCPort))
	if err != nil {
		log.Fatalf("Error when listening on TCP port: %v", err)
	}

	repo, err := redis.NewRedisTokenRepository(config.RedisAddress, config.RedisPassword)
	if err != nil {
		log.Fatal("could not connect to redis")
	}

	s, err := handlers.NewTokenService(repo,
		config.RefreshTokenSecret,
		config.RefreshDuration,
		config.AccessDuration,
	)
	if err != nil {
		log.Fatalf("Error creating token service: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTokenServiceServer(grpcServer, s)

	errChan := make(chan error)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() { errChan <- grpcServer.Serve(lis) }()

	log.Println("Starting token service...")
	log.Printf("Listening on port: %s", config.GRPCPort)

	select {
	case <-quit:
		grpcServer.GracefulStop()
	case err := <-errChan:
		log.Fatalf("GRPC Server error: %v", err)
	}

}
