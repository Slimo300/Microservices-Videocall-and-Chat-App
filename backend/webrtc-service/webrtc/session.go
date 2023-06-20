package webrtc

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

const PING_INTERVAL = 50 * time.Second

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectRoom(r *Room, w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("Upgrader error: %v\n", err)
		return
	}

	ws := &threadSafeWriter{Conn: conn, Mutex: sync.Mutex{}}
	defer ws.Close()

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
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

	r.ListLock.Lock()
	r.PeerConnections = append(r.PeerConnections, peerConnectionState{
		peerConnection: peerConnection,
		websocket:      ws,
	})
	r.ListLock.Unlock()

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
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			r.SignalPeerConnections()
		}
	})

	peerConnection.OnTrack(func(tr *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		trackLocal := r.AddTrack(tr)
		defer r.RemoveTrack(trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := tr.Read(buf)
			if err != nil {
				return
			}
			if _, err = trackLocal.Write(buf[:i]); err != nil {
				return
			}
		}
	})

	r.SignalPeerConnections()

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
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				log.Printf("Error unmarshaling answer: %v\n", err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				log.Printf("Error setting remote description: %v\n", err)
				return
			}
		}
	}

}
