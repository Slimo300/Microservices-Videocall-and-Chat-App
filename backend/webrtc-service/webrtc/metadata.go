package webrtc

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/pion/webrtc/v3"
)

type MetadataSignaler struct {
	signalerLock sync.Mutex
	dataChannels map[string]*webrtc.DataChannel
	userStreams  map[string]string
}

func NewMetadataSignaler() *MetadataSignaler {
	return &MetadataSignaler{
		signalerLock: sync.Mutex{},
		dataChannels: make(map[string]*webrtc.DataChannel),
		userStreams:  make(map[string]string),
	}
}

func (ms *MetadataSignaler) AddChannel(username string, channel *webrtc.DataChannel) {
	ms.signalerLock.Lock()
	defer ms.signalerLock.Unlock()

	channel.OnClose(func() {
		ms.signalerLock.Lock()
		defer ms.signalerLock.Unlock()

		delete(ms.dataChannels, username)
		delete(ms.userStreams, username)
	})

	channel.OnMessage(func(msg webrtc.DataChannelMessage) {
		var message dataChannelMessage
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshalling dataChannelMessage: %v", err)
			return
		}

		switch message.Type {

		case "NewUser":
			username, userOk := message.Data["username"]
			streamID, streamOk := message.Data["streamID"]

			if !userOk || !streamOk {
				log.Printf("Error unmarshaling message")
				return
			}

			ms.HandleNewUser(channel, newUserMessage{Username: username, StreamID: streamID}, msg.Data)

		case "TrackModified":
			streamID, streamOk := message.Data["streamID"]
			trackID, trackOk := message.Data["trackID"]
			isMuted, isMutedOk := message.Data["isActive"]

			if !streamOk || !trackOk || !isMutedOk {
				log.Printf("Error unmarshaling message")
				return
			}

			ms.HandleTrackModified(trackModifiedMessage{StreamID: streamID, TrackID: trackID, IsMuted: isMuted}, msg.Data)
		}
	})

	ms.dataChannels[username] = channel
}

func (ms *MetadataSignaler) HandleNewUser(newUserChannel *webrtc.DataChannel, newUserMsg newUserMessage, originalMsg []byte) {
	ms.signalerLock.Lock()
	defer ms.signalerLock.Unlock()

	// send to all data channel info about new user, also send it to him
	for username, dc := range ms.dataChannels {
		if username == newUserMsg.Username {
			continue
		}
		if err := dc.SendText(string(originalMsg)); err != nil {
			log.Printf("Error sending msg: %v", err)
		}
	}

	// send to new user info about users already in session
	for username, streamID := range ms.userStreams {

		msg, err := json.Marshal(dataChannelMessage{
			Type: "NewUser",
			Data: map[string]string{
				"username": username,
				"streamID": streamID,
			},
		})
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			return
		}

		if err := newUserChannel.SendText(string(msg)); err != nil {
			log.Printf("Couldn't send message: %v", err)
			return
		}
	}

	// save StreamID in cache
	ms.userStreams[newUserMsg.Username] = newUserMsg.StreamID
}

func (ms *MetadataSignaler) HandleTrackModified(trackMessage trackModifiedMessage, originalMessage []byte) {
	ms.signalerLock.Lock()
	defer ms.signalerLock.Unlock()

	var modifier string

	for username, stream := range ms.userStreams {
		if stream == trackMessage.StreamID {
			modifier = username
		}
	}

	for username, dc := range ms.dataChannels {
		if username != modifier {
			if err := dc.SendText(string(originalMessage)); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}

type dataChannelMessage struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

type newUserMessage struct {
	Username string `json:"username"`
	StreamID string `json:"streamID"`
}

type trackModifiedMessage struct {
	StreamID string `json:"streamID"`
	TrackID  string `json:"trackID"`
	IsMuted  string `json:"isMuted"`
}
