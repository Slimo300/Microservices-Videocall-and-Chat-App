package webrtc

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

type Room struct {
	// lock for peerConnections and trackLocals
	ListLock    sync.RWMutex
	Clients     []client
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
	closeChan   chan struct{}
	closed      bool
}

func NewRoom() *Room {
	room := &Room{}
	room.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room.closeChan = make(chan struct{})
	room.closed = false

	return room
}

func (r *Room) Done() <-chan struct{} {
	return r.closeChan
}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type client struct {
	peerConnection *webrtc.PeerConnection
	websocket      *threadSafeWriter
	userData       UserConnData
}

func (r *Room) AddClient(peerConnection *webrtc.PeerConnection, ws *threadSafeWriter, userData UserConnData) {

	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.SignalPeerConnections()
	}()

	data, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error marshaling userData: %v", err)
	}

	for _, client := range r.Clients {
		// send client info about new user
		client.websocket.WriteJSON(&websocketMessage{
			Event: "user_info",
			Data:  string(data),
		})

		clientData, err := json.Marshal(client.userData)
		if err != nil {
			log.Printf("Error marshaling userData: %v", err)
		}

		// send new user info about client
		ws.WriteJSON(&websocketMessage{
			Event: "user_info",
			Data:  string(clientData),
		})
	}

	r.Clients = append(r.Clients, client{
		peerConnection: peerConnection,
		websocket:      ws,
		userData:       userData,
	})

}

func (r *Room) ToggleMute(streamID string, videoEnabled, audioEnabled *bool) {
	r.ListLock.Lock()
	defer r.ListLock.Unlock()

	data, err := json.Marshal(UserConnData{
		StreamID:     streamID,
		VideoEnabled: videoEnabled,
		AudioEnabled: audioEnabled,
	})
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
	}

	for i := range r.Clients {
		if r.Clients[i].userData.StreamID == streamID {
			if videoEnabled != nil {
				r.Clients[i].userData.VideoEnabled = videoEnabled
			}
			if audioEnabled != nil {
				r.Clients[i].userData.AudioEnabled = audioEnabled
			}
		} else {
			r.Clients[i].websocket.WriteJSON(&websocketMessage{
				Event: "mute",
				Data:  string(data),
			})
		}
	}
}

// Add to list of tracks and fire renegotation for all PeerConnections
func (r *Room) AddTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.SignalPeerConnections()
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	r.TrackLocals[t.ID()] = trackLocal
	return trackLocal
}

// Remove from list of tracks and fire renegotation for all PeerConnections
func (r *Room) RemoveTrack(t *webrtc.TrackLocalStaticRTP) {
	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.SignalPeerConnections()
	}()

	delete(r.TrackLocals, t.ID())
}

// signalPeerConnections updates each PeerConnection so that it is getting all the expected media tracks
func (r *Room) SignalPeerConnections() {
	log.Println("Syncing connections")
	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.DispatchKeyFrame()
	}()

	attemptSync := func() (tryAgain bool) {

		if len(r.Clients) == 0 && len(r.TrackLocals) == 0 {
			log.Printf("No clients in room")
			if !r.closed {
				close(r.closeChan)
				r.closed = true
			}
			return false
		}

		for i := range r.Clients {
			if r.Clients[i].peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				r.Clients = append(r.Clients[:i], r.Clients[i+1:]...)
				log.Printf("Signal: Removing client")
				return true // We modified the slice, start from the beginning
			}

			// map of sender we already are sending, so we don't double send
			existingSenders := map[string]bool{}

			for _, sender := range r.Clients[i].peerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				// If we have a RTPSender that doesn't map to a existing track remove and signal
				if _, ok := r.TrackLocals[sender.Track().ID()]; !ok {
					if err := r.Clients[i].peerConnection.RemoveTrack(sender); err != nil {
						log.Printf("Signal: client has deleted track. Removing it...")
						return true
					}
				}
			}

			// Don't receive videos we are sending, make sure we don't have loopback
			for _, receiver := range r.Clients[i].peerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			// Add all track we aren't sending yet to the PeerConnection
			for trackID := range r.TrackLocals {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := r.Clients[i].peerConnection.AddTrack(r.TrackLocals[trackID]); err != nil {
						log.Printf("Client doesn't have a track that is currently tracked. Adding it...")
						return true
					}
				}
			}

			offer, err := r.Clients[i].peerConnection.CreateOffer(nil)
			if err != nil {
				log.Printf("Signal: Error creating offer: %v", err)
				return true
			}

			if err = r.Clients[i].peerConnection.SetLocalDescription(offer); err != nil {
				log.Printf("Signal: Error setting local description: %v", err)
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				log.Printf("Signal: Error marshaling offer: %v", err)
				return true
			}

			if err = r.Clients[i].websocket.WriteJSON(&websocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				log.Printf("Signal: Error writing to websocket: %v", err)
				return true
			}
		}

		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				r.SignalPeerConnections()
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}

// dispatchKeyFrame sends a keyframe to all PeerConnections, used everytime a new user joins the call
func (r *Room) DispatchKeyFrame() {
	r.ListLock.Lock()
	defer r.ListLock.Unlock()

	for i := range r.Clients {
		for _, receiver := range r.Clients[i].peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			_ = r.Clients[i].peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
