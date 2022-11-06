package handlers

import (
	"fmt"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/communication"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/ws"
)

type Server struct {
	DB           database.DBLayer
	Hub          ws.Hub
	TokenService auth.TokenClient
	actionChan   chan<- *communication.Action
	messageChan  <-chan *communication.Message
}

func (s *Server) RunHub() {
	go s.ListenToHub()
	s.Hub.Run()
}

func (s *Server) ListenToHub() {
	var msg *communication.Message
	for {
		select {
		case msg = <-s.messageChan:
			when, err := time.Parse(communication.TIME_FORMAT, msg.When)
			if err != nil {
				panic(err.Error())
			}
			fmt.Print(when)
			// send message to kafka/rabbit
		}
	}
}
