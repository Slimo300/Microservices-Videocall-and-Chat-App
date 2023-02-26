package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth/pb"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/configuration"
	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/repo/redis"
	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/server"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
)

func main() {

	config, err := configuration.LoadConfig(os.Getenv("CHAT_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.TokenService.GRPCPort))
	if err != nil {
		log.Fatalf("Error when listening on TCP port: %v", err)
	}

	priv, err := ioutil.ReadFile(config.TokenService.AccessTokenPrivateKey)
	if err != nil {
		log.Fatalf("could not read private key pem file: %v", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		log.Fatalf("could not parse private key: %v", err)
	}

	repo, err := redis.NewRedisTokenRepository(config.TokenService.RedisAddress, config.TokenService.RedisPass)
	if err != nil {
		log.Fatal("could not connect to redis")
	}

	s := server.NewTokenService(repo,
		config.TokenService.RefreshTokenSecret,
		*privKey,
		config.RefreshDuration,
		config.AccessDuration,
	)

	grpcServer := grpc.NewServer()

	pb.RegisterTokenServiceServer(grpcServer, s)

	errChan := make(chan error)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() { errChan <- grpcServer.Serve(lis) }()

	log.Println("Starting token service...")
	log.Printf("Listening on port: %s", config.TokenService.GRPCPort)

	select {
	case <-quit:
		grpcServer.GracefulStop()
	case err := <-errChan:
		log.Fatalf("GRPC Server error: %v", err)
	}

}
