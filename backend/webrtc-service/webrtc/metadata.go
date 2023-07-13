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
		}

		switch message.Type {
		case "NewUser":
			var messageBody newUserMessage
			if err := json.Unmarshal(message.Data, &messageBody); err != nil {
				log.Printf("Error unmarshaling newUserMessage: %v", err)
			}

			ms.HandleNewUser(channel, messageBody, msg.Data)

		case "TrackModified":
			var messageBody trackModifiedMessage
			if err := json.Unmarshal(message.Data, &messageBody); err != nil {
				log.Printf("Error unmarshaling newUserMessage: %v", err)
			}

			ms.HandleTrackModified(messageBody, msg.Data)
		}
	})

	ms.dataChannels[username] = channel
}

func (ms *MetadataSignaler) HandleNewUser(newUserChannel *webrtc.DataChannel, newUserMsg newUserMessage, originalMsg []byte) {
	ms.signalerLock.Lock()
	defer ms.signalerLock.Unlock()

	// send to all data channel info about new user, also send it to him
	for _, dc := range ms.dataChannels {
		dc.SendText(string(originalMsg))
	}

	// send to new user info about users already in session
	for username, streamID := range ms.userStreams {
		msgBody, err := json.Marshal(newUserMessage{
			Username: username,
			StreamID: streamID,
		})
		if err != nil {
			log.Printf("Error marshaling msgBody: %v", err)
		}

		msg, err := json.Marshal(dataChannelMessage{
			Type: "NewUser",
			Data: msgBody,
		})

		if err := newUserChannel.SendText(string(msg)); err != nil {
			log.Printf("Couldn't send message: %v", err)
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
			dc.SendText(string(originalMessage))
		}
	}
}

type dataChannelMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type newUserMessage struct {
	Username string `json:"username"`
	StreamID string `json:"streamID"`
}

type trackModifiedMessage struct {
	StreamID string `json:"streamID"`
	TrackID  string `json:"trackID"`
	IsMuted  bool   `json:"isMuted"`
}
