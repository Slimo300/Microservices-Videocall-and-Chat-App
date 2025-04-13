package handlers

import (
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
)

type message struct {
	ID       string        `json:"ID"`
	GroupID  string        `json:"groupID"`
	Posted   time.Time     `json:"created"`
	Text     string        `json:"text"`
	MemberID string        `json:"memberID"`
	Member   member        `json:"member"`
	Files    []messageFile `json:"files"`
}

type member struct {
	ID       string `json:"ID"`
	GroupID  string `json:"groupID"`
	UserID   string `json:"userID"`
	Username string `json:"username"`
}

type messageFile struct {
	Key       string `json:"key"`
	Extension string `json:"ext"`
}

func modelMessageToResponse(msg models.Message) message {
	files := []messageFile{}
	for _, f := range msg.Files() {
		files = append(files, modelMessageFileToResponse(f))
	}

	return message{
		ID:       msg.ID().String(),
		GroupID:  msg.GroupID().String(),
		Posted:   msg.Posted(),
		Text:     msg.Text(),
		MemberID: msg.MemberID().String(),
		Member:   modelMemberToResponse(msg.Member()),
		Files:    files,
	}
}

func modelMemberToResponse(mem models.Member) member {
	return member{
		ID:       mem.ID().String(),
		GroupID:  mem.GroupID().String(),
		UserID:   mem.UserID().String(),
		Username: mem.Username(),
	}
}

func modelMessageFileToResponse(file models.MessageFile) messageFile {
	return messageFile{
		Key:       file.Key(),
		Extension: file.Extension(),
	}
}

func modelMessagesToResponse(messageModels []models.Message) []message {
	msgs := []message{}
	for _, msg := range messageModels {
		msgs = append(msgs, modelMessageToResponse(msg))
	}
	return msgs
}
