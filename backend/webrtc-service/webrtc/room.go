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
	ListLock   sync.RWMutex
	turnConfig webrtc.Configuration

	Peers       []*Peer
	TrackLocals map[string]*LocalTrack

	BanningRules map[MutingRule]bool
	MutingRules  map[MutingRule]bool

	PeerRights map[string]PeerRights

	// closeChan is a channel on which room signals that it has exited, closed flag makes sure closeChan is only closed once
	closeChan chan struct{}
	closed    bool
}

func NewRoom(config webrtc.Configuration) *Room {
	room := &Room{}
	room.TrackLocals = make(map[string]*LocalTrack)

	room.MutingRules = make(map[MutingRule]bool)
	room.BanningRules = make(map[MutingRule]bool)
	room.PeerRights = make(map[string]PeerRights)

	room.closeChan = make(chan struct{})
	room.closed = false
	room.turnConfig = config

	return room
}

type LocalTrack struct {
	Track    *webrtc.TrackLocalStaticRTP
	MemberID string
}

func (r *Room) Done() <-chan struct{} {
	return r.closeChan
}

func (r *Room) GetPeerRights(memberID string) (PeerRights, bool) {
	r.ListLock.RLock()
	defer r.ListLock.RUnlock()

	rights, ok := r.PeerRights[memberID]

	return rights, ok
}

func (r *Room) AddPeer(peer *Peer) {
	log.Println("ROOM: Adding Peer")

	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.SignalPeerConnections()
	}()

	data, err := json.Marshal(peer.userData)
	if err != nil {
		log.Printf("Error marshaling userData: %v", err)
	}

	r.PeerRights[peer.userData.MemberID] = PeerRights{
		Creator: peer.userData.Creator,
		Admin:   peer.userData.Admin,
		Muting:  peer.userData.Muting,
	}

	for _, client := range r.Peers {
		// send current peers info about new peer
		if err := client.signaler.WriteJSON(&websocketMessage{
			Event: "user_info",
			Data:  string(data),
		}); err != nil {
			log.Println(err.Error())
		}

		clientData, err := json.Marshal(client.userData)
		if err != nil {
			log.Default().Printf("Error marshaling userData: %v", err)
		}

		// send new peer info about current peers
		if err := peer.signaler.WriteJSON(&websocketMessage{
			Event: "user_info",
			Data:  string(clientData),
		}); err != nil {
			log.Println(err.Error())
		}
	}

	r.Peers = append(r.Peers, peer)
}

func (r *Room) AddMutingRule(mr MutingRule) {
	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.SignalPeerConnections()
	}()

	r.MutingRules[mr] = true
}

func (r *Room) RemoveMutingRule(mr MutingRule) {
	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.SignalPeerConnections()
	}()

	delete(r.MutingRules, mr)
}
func (r *Room) AddBanningRule(mr MutingRule) {
	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.SignalPeerConnections()
	}()

	r.BanningRules[mr] = true
}

func (r *Room) RemoveBanningRule(mr MutingRule) {
	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.SignalPeerConnections()
	}()

	delete(r.BanningRules, mr)
}

func (r *Room) SignalPeerClosed(memberID string) {
	r.ListLock.RLock()
	defer r.ListLock.RUnlock()

	for i := range r.Peers {
		if r.Peers[i].userData.MemberID != memberID {
			r.Peers[i].signaler.WriteJSON(&websocketMessage{
				Event: "disconnected",
				Data:  memberID,
			})
		}
	}
}

// Add to list of tracks and fire renegotation for all PeerConnections
func (r *Room) AddTrack(t *webrtc.TrackRemote, memberID string) *webrtc.TrackLocalStaticRTP {
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

	r.TrackLocals[t.ID()] = &LocalTrack{
		Track:    trackLocal,
		MemberID: memberID,
	}
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
	r.ListLock.Lock()
	defer func() {
		r.ListLock.Unlock()
		r.DispatchKeyFrame()
	}()

	attemptSync := func() bool {
		if len(r.Peers) == 0 && len(r.TrackLocals) == 0 {
			if !r.closed {
				close(r.closeChan)
				r.closed = true
			}
			return false
		}

		for i := range r.Peers {
			if r.Peers[i].peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				r.Peers = append(r.Peers[:i], r.Peers[i+1:]...)
				return true // We modified the slice, start from the beginning
			}

			// createOffer flag shows whether a renegotiation is needed
			// map of sender we already are sending, so we don't double send
			existingSenders := map[string]bool{}

			for _, sender := range r.Peers[i].peerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				track, ok := r.TrackLocals[sender.Track().ID()]
				// remove track if its not in tracks tracked by room or if there is a muting rule associated with it
				if !ok {
					if err := r.Peers[i].peerConnection.RemoveTrack(sender); err != nil {
						log.Printf("Signal: Couldn't remove track from peerConnection: %v", err)
						return true
					}
				} else {
					for mr := range r.MutingRules {
						log.Println(mr.MemberID)
					}
					mr := MutingRule{MemberID: track.MemberID, TrackKind: track.Track.Kind().String()}
					if r.MutingRules[mr] || r.BanningRules[mr] || r.Peers[i].mutingRules[mr] {
						log.Printf("Track removed from peerConnection because its muted")
						if err := r.Peers[i].peerConnection.RemoveTrack(sender); err != nil {
							log.Printf("Signal: Couldn't remove track from peerConnection: %v", err)
							return true
						}
					}
				}
			}

			// Don't receive videos we are sending, make sure we don't have loopback
			for _, receiver := range r.Peers[i].peerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}
				existingSenders[receiver.Track().ID()] = true
			}

			// Add tracks we aren't sending yet to the PeerConnection
			for trackID, track := range r.TrackLocals {
				mr := MutingRule{MemberID: track.MemberID, TrackKind: track.Track.Kind().String()}
				_, ok := existingSenders[trackID]
				// Add track if its not in existing senders and if there are no mutingRules disallowing it
				if !ok && !r.MutingRules[mr] && !r.BanningRules[mr] && !r.Peers[i].mutingRules[mr] {
					if _, err := r.Peers[i].peerConnection.AddTrack(r.TrackLocals[trackID].Track); err != nil {
						log.Printf("Error adding track to peerConnection: %v", err)
						return true
					}
				}
			}

			offer, err := r.Peers[i].peerConnection.CreateOffer(nil)
			if err != nil {
				log.Printf("Signal: Error creating offer: %v", err)
				return true
			}

			if err = r.Peers[i].peerConnection.SetLocalDescription(offer); err != nil {
				log.Printf("Signal: Error setting local description: %v", err)
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				log.Printf("Signal: Error marshaling offer: %v", err)
				return true
			}

			if err = r.Peers[i].signaler.WriteJSON(&websocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				log.Printf("Signal: Error writing to websocket: %v", err)
				return true
			}
		}
		return false
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

	for i := range r.Peers {
		for _, receiver := range r.Peers[i].peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			_ = r.Peers[i].peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
