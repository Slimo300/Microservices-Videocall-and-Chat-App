package webrtc

import (
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
)

type RoomsRelay struct {
	rooms      map[string]*Room
	relayMutex sync.Mutex
	turnConfig webrtc.Configuration
}

func NewRoomsRelay(turnConfig webrtc.Configuration) *RoomsRelay {
	return &RoomsRelay{
		rooms:      map[string]*Room{},
		relayMutex: sync.Mutex{},
		turnConfig: turnConfig,
	}
}

func (r *RoomsRelay) GetRoom(groupID string) *Room {
	r.relayMutex.Lock()
	defer r.relayMutex.Unlock()

	room, ok := r.rooms[groupID]
	if !ok {
		room = NewRoom(r.turnConfig)
		r.rooms[groupID] = room
	}

	go func() {
		for {
			select {
			case <-time.NewTicker(3 * time.Second).C:
				room.DispatchKeyFrame()
			case <-room.Done():
				r.relayMutex.Lock()
				delete(r.rooms, groupID)
				r.relayMutex.Unlock()
				return
			}
		}
	}()

	return room
}
