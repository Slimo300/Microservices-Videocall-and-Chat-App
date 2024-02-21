package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	rtc "github.com/pion/webrtc/v3"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/builder"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/database/redis"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/eventprocessor"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/webrtc"
)

func getPublicKey() (*rsa.PublicKey, error) {

	bytePubKey, err := os.ReadFile("/rsa/public.key")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytePubKey)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func main() {

	conf, err := config.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	pubKey, err := getPublicKey()
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	db, err := redis.Setup(conf.DBAddress, conf.DBPassword)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	builder, err := builder.NewBrokerBuilder(msgqueue.ParseBrokerType(conf.BrokerType), conf.BrokerAddress)
	if err != nil {
		log.Fatalf("Error creating broker builder: %v", err)
	}

	emitter, err := builder.GetEmiter(msgqueue.EmiterConfig{
		ExchangeName: "webrtc",
	})
	if err != nil {
		log.Fatalf("Error building emitter: %v", err)
	}

	listener, err := builder.GetListener(msgqueue.ListenerConfig{
		ClientName: "webrtc-service",
		Events: []msgqueue.Event{
			events.GroupDeletedEvent{},
			events.MemberCreatedEvent{},
			events.MemberUpdatedEvent{},
			events.MemberDeletedEvent{},
		},
	})
	if err != nil {
		log.Fatalf("Error building listener: %v", err)
	}

	if err := emitter.Emit(events.ServiceStartedEvent{
		ServiceAddress: fmt.Sprintf("%s.%s.%s.svc.cluster.local:%s", conf.PodName, conf.ServiceName, conf.PodNamespace, conf.HTTPPort),
	}); err != nil {
		log.Fatalf("Couldn't emit ServiceStartedEvent")
	}

	go eventprocessor.NewDBEventProcessor(listener, db).ProcessEvents("group")

	turnConfig := rtc.Configuration{
		ICEServers: []rtc.ICEServer{
			{
				URLs: []string{fmt.Sprintf("stun:%s:%s", conf.TURNAddress, conf.TURNPort)},
			},
			{
				URLs:           []string{fmt.Sprintf("turn:%s:%s", conf.TURNAddress, conf.TURNPort)},
				Username:       conf.TURNUser,
				Credential:     conf.TURNPassword,
				CredentialType: rtc.ICECredentialTypePassword,
			},
			{
				URLs:           []string{fmt.Sprintf("turns:%s:%s", conf.TURNAddress, conf.TURNSPort)},
				Username:       conf.TURNUser,
				Credential:     conf.TURNPassword,
				CredentialType: rtc.ICECredentialTypePassword,
			},
			{
				URLs:           []string{fmt.Sprintf("turns:%s:%s?transport=tcp", conf.TURNAddress, conf.TURNSPort)},
				Username:       conf.TURNUser,
				Credential:     conf.TURNPassword,
				CredentialType: rtc.ICECredentialTypePassword,
			},
		},
	}

	server := &handlers.Server{
		DB:        db,
		PublicKey: pubKey,
		Relay:     webrtc.NewRoomsRelay(turnConfig),
	}

	handler := server.Setup(conf.Origin)

	httpServer := http.Server{
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
