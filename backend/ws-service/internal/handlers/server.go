package handlers

import (
	"log"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue"
	tokens "github.com/Slimo300/chat-tokenservice/pkg/client"
	"github.com/Slimo300/chat-wsservice/internal/cache"
	"github.com/Slimo300/chat-wsservice/internal/database"
	"github.com/Slimo300/chat-wsservice/internal/ws"
)

type Server struct {
	DB          database.DBLayer
	Hub         *ws.WSHub
	TokenClient tokens.TokenClient
	CodeCache   cache.AccessCodeCache
	Emitter     msgqueue.EventEmiter
	Listener    msgqueue.EventListener
	MessageChan <-chan *ws.Message
	EventChan   chan<- msgqueue.Event
}

func (s *Server) RunHub() {
	go s.ListenToHub()
	s.Hub.Run()
}

func (s *Server) ListenToHub() {
	for msg := range s.MessageChan {
		var files []events.File
		for _, f := range msg.Files {
			files = append(files, events.File{Key: f.Key, Extension: f.Ext})
		}

		if err := s.Emitter.Emit(events.MessageSentEvent{
			ID:      msg.ID,
			GroupID: msg.Group,
			UserID:  msg.User,
			Nick:    msg.Nick,
			Posted:  msg.When,
			Text:    msg.Message,
			Files:   files,
		}); err != nil {
			log.Printf("Error when sending message to broker %v: %v", msg, err)
		}
	}
}
