package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/config"
	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/database/redis"
	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/handlers"
	"google.golang.org/grpc"
)

func readPrivateKey() (*rsa.PrivateKey, error) {

	bytePrivKey, err := os.ReadFile("/rsa/private.key")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytePrivKey)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil

}

func main() {

	privateKey, err := readPrivateKey()
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}

	config, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Couldn't read configuration: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCPort))
	if err != nil {
		log.Fatalf("Error when listening on TCP port: %v", err)
	}

	db, err := redis.NewRedisTokenDB(config.RedisAddress, config.RedisPassword)
	if err != nil {
		log.Fatal("could not connect to redis")
	}

	s, err := handlers.NewTokenService(db,
		privateKey,
		config.RefreshTokenSecret,
		config.RefreshDuration,
		config.AccessDuration,
	)
	if err != nil {
		log.Fatalf("Error creating token service: %v", err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterTokenServiceServer(grpcServer, s)

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
