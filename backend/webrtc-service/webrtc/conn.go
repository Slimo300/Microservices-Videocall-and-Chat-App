package webrtc

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type UserConnData struct {
	Username     string `json:"username,omitempty"`
	StreamID     string `json:"streamID,omitempty"`
	VideoEnabled *bool  `json:"videoEnabled,omitempty"`
	AudioEnabled *bool  `json:"audioEnabled,omitempty"`
}

var (
	turnConfig = webrtc.Configuration{
		ICETransportPolicy: webrtc.ICETransportPolicyRelay,
		ICEServers: []webrtc.ICEServer{
			{

				URLs: []string{"stun:turn-around.pro:3478"},
			},
			{

				URLs: []string{"turn:turn-around.pro:3478"},

				Username: "test",

				Credential:     "test123",
				CredentialType: webrtc.ICECredentialTypePassword,
			},
		},
	}
)

func (r *Room) ConnectRoom(conn *websocket.Conn, userData UserConnData) {

	ws := newThreadSafeWriter(conn)
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

	r.AddClient(peerConnection, ws, userData)

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
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
		switch pcs {
		case webrtc.PeerConnectionStateFailed:
			log.Println("Peer connection failed")
			if err := peerConnection.Close(); err != nil {
				log.Printf("Error closing failed connection: %v", err)
			}
		case webrtc.PeerConnectionStateClosed:
			log.Printf("Peer connection closed")
			r.SignalPeerConnections()
		}
	})

	peerConnection.OnTrack(func(tr *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		log.Println("New track received!")
		trackLocal := r.AddTrack(tr)
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

	message := &websocketMessage{}
	for {
		_, raw, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Printf("Error unmarshaling message: %v\n", err)
			return
		}

		switch message.Event {
		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				log.Printf("Error unmarshaling candidate: %v\n", err)
				return
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Printf("Error adding ICE candidate: %v\n", err)
				return
			}
		case "answer":
			log.Printf("SDP Answer received")
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				log.Printf("Error unmarshaling answer: %v\n", err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				log.Printf("Error setting remote description: %v\n", err)
				return
			}
		case "renegotiate":
			r.SignalPeerConnections()
		case "mute":
			muteInfo := struct {
				VideoEnabled *bool `json:"videoEnabled,omitempty"`
				AudioEnabled *bool `json:"audioEnabled,omitempty"`
			}{}

			if err := json.Unmarshal([]byte(message.Data), &muteInfo); err != nil {
				log.Printf("Error unmarshaling mute message")
			}

			r.ToggleMute(userData.StreamID, muteInfo.VideoEnabled, muteInfo.AudioEnabled)
		}
	}

}
