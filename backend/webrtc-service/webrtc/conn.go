package webrtc

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type UserConnData struct {
	MemberID   string `json:"memberID,omitempty"`
	StreamID   string `json:"streamID,omitempty"`
	Username   string `json:"username,omitempty"`
	PictureURL string `json:"pictureURL,omitempty"`
	Muting     bool   `json:"muting,omitempty"`
	Admin      bool   `json:"admin,omitempty"`
	Creator    bool   `json:"creator,omitempty"`
}

func (r *Room) ConnectRoom(conn *websocket.Conn, userData UserConnData) {

	defer log.Println("Closing connection")
	defer r.SignalPeerClosed(userData.MemberID)
	ws := NewSignaler(conn)
	defer ws.Close()

	peerConnection, err := webrtc.NewPeerConnection(r.turnConfig)
	if err != nil {
		log.Printf("Creating peer connection error: %v\n", err)
		return
	}
	defer peerConnection.Close()

	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Printf("Error adding transceiver %d: %v", typ, err)
			return
		}
	}

	peer := NewPeer(peerConnection, ws, r, userData)
	r.AddPeer(peer)

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		log.Println("New candidate received")
		if i == nil {
			return
		}

		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Printf("Error marshaling ice candidate %v\n", err)
			return
		}

		if err := ws.WriteJSON(&websocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); err != nil {
			log.Printf("Error sending candidate through websocket: %v", err)
		}
	})

	peerConnection.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		log.Println("Connection state changed: ", pcs.String())
		switch pcs {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Printf("Error closing failed connection: %v", err)
			}
		case webrtc.PeerConnectionStateClosed:
			r.SignalPeerConnections()
		}
	})

	peerConnection.OnTrack(func(tr *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		log.Println("New Track received...")
		trackLocal := r.AddTrack(tr, userData.MemberID)
		defer r.RemoveTrack(trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := tr.Read(buf)
			if err != nil {
				log.Println("track remote read err")
				return
			}
			if _, err = trackLocal.Write(buf[:i]); err != nil {
				log.Println("track local write err")
				return
			}
		}
	})

	peer.ServeSignaler()
}
