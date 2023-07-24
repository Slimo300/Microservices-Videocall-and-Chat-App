package webrtc

import (
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
)

type RoomsRelay struct {
	rooms      map[string]*Room
	relayMutex sync.Mutex
}

func NewRoomsRelay() *RoomsRelay {
	return &RoomsRelay{
		rooms:      map[string]*Room{},
		relayMutex: sync.Mutex{},
	}
}

func (r *RoomsRelay) GetRoom(groupID string) *Room {
	r.relayMutex.Lock()
	defer r.relayMutex.Unlock()

	room, ok := r.rooms[groupID]
	if !ok {
		room = &Room{}
		room.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
		// room.DataHandler = NewMetadataSignaler()
		r.rooms[groupID] = room
	}

	go func() {
		for range time.NewTicker(3 * time.Second).C {
			room.DispatchKeyFrame()
		}
	}()

	return room
}
