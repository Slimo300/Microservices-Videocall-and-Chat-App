package command

import (
	"context"
	"sync"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"github.com/google/uuid"
)

type DeleteMessageForEveryoneCommand struct {
	MessageID, UserID uuid.UUID
}

type DeleteMessageForEveryoneHandler struct {
	repo    database.MessagesRepository
	storage storage.Storage
	emitter msgqueue.EventEmiter
}

func NewDeleteMessageForEveryoneHandler(repo database.MessagesRepository, storage storage.Storage, emitter msgqueue.EventEmiter) DeleteMessageForEveryoneHandler {
	if repo == nil {
		panic("nil repo")
	}
	if storage == nil {
		panic("nil storage")
	}
	if emitter == nil {
		panic("nil emitter")
	}
	return DeleteMessageForEveryoneHandler{repo: repo, storage: storage, emitter: emitter}
}

func (h *DeleteMessageForEveryoneHandler) Handle(ctx context.Context, cmd DeleteMessageForEveryoneCommand) error {
	msg, err := h.repo.GetMessageByID(ctx, cmd.UserID, cmd.MessageID)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	errors := make(chan error, 1)
	defer close(errors)
	for _, file := range msg.Files() {
		file := file
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.storage.DeleteFile(ctx, file.Key()); err != nil {
				select {
				case errors <- err:
				default:
				}
			}
		}()
		wg.Wait()
		if err := <-errors; err != nil {
			return err
		}
	}
	if err := h.repo.DeleteMessageForEveryone(ctx, cmd.UserID, cmd.MessageID); err != nil {
		return err
	}
	return h.emitter.Emit(events.MessageDeletedEvent{ID: msg.ID(), GroupID: msg.GroupID()})
}
