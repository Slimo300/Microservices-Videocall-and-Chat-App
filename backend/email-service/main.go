package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/email-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/email-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/email"
	"google.golang.org/grpc"
)

func main() {
	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Error when reading configuration: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPCPort))
	if err != nil {
		log.Fatalf("Error creating TCP listener: %v", err)
	}

	s, err := handlers.NewEmailService(conf.EmailFrom,
		conf.SMTPHost,
		conf.SMTPPort,
		conf.SMTPUser,
		conf.SMTPPass,
		conf.Origin,
	)
	if err != nil {
		log.Fatalf("Error when creating email service: %v", err)
	}

	grpcServer := grpc.NewServer()
	email.RegisterEmailServiceServer(grpcServer, s)

	errChan := make(chan error)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() { errChan <- grpcServer.Serve(lis) }()

	log.Println("Starting email service..")
	log.Printf("Listening on port %s", conf.GRPCPort)

	select {
	case <-quit:
		grpcServer.GracefulStop()
	case err := <-errChan:
		log.Fatalf("GRPC Server error: %v", err)
	}
}
