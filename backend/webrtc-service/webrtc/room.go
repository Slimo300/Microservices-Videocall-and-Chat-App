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

	data, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error marshaling userData: %v", err)
	}

	for _, client := range r.Clients {
		client.websocket.WriteJSON(&websocketMessage{
			Event: "newUser",
			Data:  string(data),
		})

		clientData, err := json.Marshal(client.userData)
		if err != nil {
			log.Printf("Error marshaling userData: %v", err)
		}

		ws.WriteJSON(&websocketMessage{
			Event: "newUser",
			Data:  string(clientData),
		})
	}

	r.Clients = append(r.Clients, client{
		peerConnection: peerConnection,
		websocket:      ws,
		userData:       userData,
	})

	r.ListLock.Unlock()
}

func (r *Room) ToggleVideoMute(streamID string, videoEnabled bool) {
	r.ListLock.Lock()
	defer r.ListLock.Unlock()

	data, err := json.Marshal(UserConnData{
		StreamID:     streamID,
		VideoEnabled: &videoEnabled,
	})
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
	}

	for _, client := range r.Clients {
		if client.userData.StreamID == streamID {
			client.userData.VideoEnabled = &videoEnabled
		} else {
			client.websocket.WriteJSON(&websocketMessage{
				Event: "toggleVideoMute",
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
		for i := range r.Clients {
			if r.Clients[i].peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				r.Clients = append(r.Clients[:i], r.Clients[i+1:]...)
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
						return true
					}
				}
			}

			offer, err := r.Clients[i].peerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = r.Clients[i].peerConnection.SetLocalDescription(offer); err != nil {
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}

			if err = r.Clients[i].websocket.WriteJSON(&websocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
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
