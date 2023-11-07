package webrtc

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

const PING_INTERVAL = 55 * time.Second

type messageType string

const (
	Candidate       messageType = "candidate"
	Answer          messageType = "answer"
	Renegotiate     messageType = "renegotiate"
	MuteYourself    messageType = "mute_yourself"
	MuteForEveryone messageType = "mute_for_everyone"
	MuteForYourself messageType = "mute_for_yourself"
)

type mutingActionType string

const (
	Disable mutingActionType = "disable"
	Enable  mutingActionType = "enable"
)

type Peer struct {
	peerLock       sync.RWMutex
	peerConnection *webrtc.PeerConnection
	signaler       *Signaler
	userData       UserConnData
	mutingRules    map[MutingRule]bool
	room           *Room
}

func NewPeer(pc *webrtc.PeerConnection, s *Signaler, r *Room, userData UserConnData) *Peer {
	return &Peer{
		peerConnection: pc,
		signaler:       s,
		userData:       userData,
		room:           r,
		mutingRules:    make(map[MutingRule]bool),
	}
}

func (p *Peer) ServeSignaler() {
	message := &websocketMessage{}
	for {
		_, raw, err := p.signaler.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Printf("Error unmarshaling message: %v\n", err)
			return
		}

		switch message.Event {
		case Candidate:
			if err := p.HandleCandidate(message.Data); err != nil {
				log.Printf("Error handling candidate: %v", err)
			}
		case Answer:
			if err := p.HandleAnswer(message.Data); err != nil {
				log.Printf("Error handling answer: %v", err)
			}
		case Renegotiate:
			if err := p.HandleRenegotiate(); err != nil {
				log.Printf("Error handling renegotiate: %v", err)
			}
		case MuteYourself:
			if err := p.HandleMuteYourself(message.Data); err != nil {
				log.Printf("Error handling mute_yourself: %v", err)
			}
		case MuteForEveryone:
			if err := p.HandleMuteForEveryone(message.Data); err != nil {
				log.Printf("Error handling mute_for_everyone: %v", err)
			}
		case MuteForYourself:
			if err := p.HandleMuteForYourself(message.Data); err != nil {
				log.Printf("Error handling mute_for_yourself: %v", err)
			}
		default:
			log.Printf("Unsupported message event type: %s", message.Event)
		}
	}
}

func (p *Peer) HandleCandidate(data string) error {
	log.Println("Handling candidate... ")
	candidate := webrtc.ICECandidateInit{}
	if err := json.Unmarshal([]byte(data), &candidate); err != nil {
		return fmt.Errorf("Error unmarshaling candidate: %v\n", err)
	}

	if err := p.peerConnection.AddICECandidate(candidate); err != nil {
		return fmt.Errorf("Error adding ICE candidate: %v\n", err)
	}

	return nil
}
func (p *Peer) HandleAnswer(data string) error {
	log.Println("Handling answer... ")
	answer := webrtc.SessionDescription{}
	if err := json.Unmarshal([]byte(data), &answer); err != nil {
		return fmt.Errorf("Error unmarshaling answer: %v\n", err)
	}

	if err := p.peerConnection.SetRemoteDescription(answer); err != nil {
		return fmt.Errorf("Error setting remote description: %v\n", err)
	}

	return nil
}

func (p *Peer) HandleRenegotiate() error {
	log.Println("Handling renegotiate... ")
	p.room.SignalPeerConnections()
	return nil
}

type MutingAction struct {
	MutingRule
	ActionType mutingActionType `json:"actionType,required"`
}

type MutingRule struct {
	MemberID  string `json:"memberID"`
	TrackKind string `json:"kind"`
}

func (p *Peer) HandleMuteYourself(data string) error {
	log.Println("Handling muting yourself")

	var action MutingAction
	if err := json.Unmarshal([]byte(data), &action); err != nil {
		return fmt.Errorf("Error unmarshaling mute_for_everyone: %w", err)
	}

	action.MutingRule.MemberID = p.userData.MemberID

	switch action.ActionType {
	case Enable:
		p.room.RemoveMutingRule(action.MutingRule)
	case Disable:
		p.room.AddMutingRule(action.MutingRule)
	}

	return nil
}

func (p *Peer) HandleMuteForEveryone(data string) error {
	log.Println("Handling MuteForEveryone")
	var action MutingAction
	if err := json.Unmarshal([]byte(data), &action); err != nil {
		return fmt.Errorf("Error unmarshaling mute_for_everyone")
	}

	if !p.CanBan(action.MutingRule) {
		return fmt.Errorf("User not authorized to mute")
	}

	switch action.ActionType {
	case Enable:
		p.room.RemoveBanningRule(action.MutingRule)
	case Disable:
		p.room.AddBanningRule(action.MutingRule)
	}

	return nil
}

func (p *Peer) HandleMuteForYourself(data string) error {
	log.Println("Handling mute for yourself: ", data)
	var action MutingAction
	if err := json.Unmarshal([]byte(data), &action); err != nil {
		return fmt.Errorf("Error unmarshaling mute_for_everyone")
	}

	switch action.ActionType {
	case Enable:
		p.RemoveMutingRule(action.MutingRule)
	case Disable:
		p.AddMutingRule(action.MutingRule)
	}

	return nil
}

func (p *Peer) CanBan(mr MutingRule) bool {
	if p.userData.MemberID == mr.MemberID {
		return false
	}

	mutedMember, ok := p.room.GetPeerRights(mr.MemberID)
	if !ok {
		return false
	}

	if p.userData.Creator {
		return true
	}

	if p.userData.Admin && !mutedMember.Creator {
		return true
	}
	if p.userData.Muting && !mutedMember.Creator && !mutedMember.Admin {
		return true
	}

	return false
}

func (p *Peer) AddMutingRule(mr MutingRule) {
	p.peerLock.Lock()
	defer func() {
		p.peerLock.Unlock()
		p.room.SignalPeerConnections()
	}()
	log.Println("Adding muting rule: ", mr)

	p.mutingRules[mr] = true
	log.Println(p.mutingRules)
}

func (p *Peer) RemoveMutingRule(mr MutingRule) {
	p.peerLock.Lock()
	defer func() {
		p.peerLock.Unlock()
		p.room.SignalPeerConnections()
	}()
	log.Println("Removing muting rule: ", mr)

	delete(p.mutingRules, mr)
	log.Println(p.mutingRules)
}

type websocketMessage struct {
	Event messageType `json:"event"`
	Data  string      `json:"data"`
}

// Helper to make Gorilla Websockets threadsafe
type Signaler struct {
	*websocket.Conn
	sync.Mutex
	ticker    *time.Ticker
	closeChan chan struct{}
}

func NewSignaler(conn *websocket.Conn) *Signaler {

	s := &Signaler{Conn: conn, Mutex: sync.Mutex{}, ticker: time.NewTicker(PING_INTERVAL), closeChan: make(chan struct{})}

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.Lock()
				if err := s.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Error pinging websocket: %v", err)
					return
				}
				s.Unlock()
			case <-s.closeChan:
				s.ticker.Stop()
				return
			}
		}
	}()

	return s
}

func (s *Signaler) WriteJSON(v interface{}) error {
	s.Lock()
	defer s.Unlock()
	s.ticker.Reset(PING_INTERVAL)
	return s.Conn.WriteJSON(v)
}

func (s *Signaler) Close() {
	s.closeChan <- struct{}{}
	s.Conn.Close()
}

type PeerRights struct {
	Creator bool
	Admin   bool
	Muting  bool
}
