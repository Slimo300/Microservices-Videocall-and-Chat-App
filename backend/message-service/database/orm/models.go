package orm

import (
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/google/uuid"
)

type Member struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	UserID           uuid.UUID `gorm:"uniqueIndex:idx_first;size:191"`
	GroupID          uuid.UUID `gorm:"uniqueIndex:idx_first;size:191"`
	Username         string
	Creator          bool
	Admin            bool
	DeletingMessages bool
}

func (Member) TableName() string {
	return "members"
}

func unmarshalMember(member models.Member) Member {
	return Member{
		ID:               member.ID(),
		UserID:           member.UserID(),
		GroupID:          member.GroupID(),
		Username:         member.Username(),
		Creator:          member.Creator(),
		Admin:            member.Admin(),
		DeletingMessages: member.DeletingMessages(),
	}
}

type Message struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	GroupID  uuid.UUID `gorm:"index"`
	Posted   time.Time
	Text     string
	MemberID uuid.UUID `gorm:"size:191"`
	Member   Member
	Deleters []Member      `gorm:"many2many:users_who_deleted;constraint:OnDelete:CASCADE;"`
	Files    []MessageFile `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE;"`
}

func unmarshalMessage(message models.Message) Message {
	var files []MessageFile
	for _, f := range message.Files() {
		files = append(files, unmarshalMessageFile(f))
	}
	// here we set files slice to nil in order to prevent GORM from creating blank message file objects
	if len(files) == 0 {
		files = nil
	}
	var deleters []Member
	for _, d := range message.Deleters() {
		deleters = append(deleters, unmarshalMember(d))
	}
	// here we set deleters slice to nil in order to prevent GORM from creating blank Member objects
	if len(deleters) == 0 {
		deleters = nil
	}
	return Message{
		ID:       message.ID(),
		GroupID:  message.GroupID(),
		MemberID: message.MemberID(),
		Member:   unmarshalMember(message.Member()),
		Text:     message.Text(),
		Posted:   message.Posted(),
		Files:    files,
		Deleters: deleters,
	}
}

func (m Message) marshalMessage() models.Message {
	var deleters []models.Member
	for _, d := range m.Deleters {
		deleters = append(deleters, models.UnmarshalMemberFromDatabase(d.ID, d.UserID, d.GroupID, d.Username, d.Creator, d.Admin, d.DeletingMessages))
	}
	if len(deleters) == 0 {
		deleters = nil
	}
	var files []models.MessageFile
	for _, f := range m.Files {
		files = append(files, models.NewMessageFile(f.MessageID, f.Key, f.Extension))
	}
	if len(files) == 0 {
		files = nil
	}
	return models.UnmarshalMessageFromDatabase(m.ID, m.GroupID, m.MemberID, m.Text, m.Posted,
		models.UnmarshalMemberFromDatabase(m.MemberID, m.Member.UserID, m.Member.GroupID, m.Member.Username, m.Member.Creator, m.Member.Admin, m.Member.DeletingMessages),
		deleters, files,
	)
}

func (Message) TableName() string {
	return "messages"
}

type MessageFile struct {
	Key       string    `gorm:"primaryKey"`
	MessageID uuid.UUID `gorm:"size:191"`
	Extension string
}

func unmarshalMessageFile(file models.MessageFile) MessageFile {
	return MessageFile{
		MessageID: file.MessageID(),
		Key:       file.Key(),
		Extension: file.Extension(),
	}
}

func (MessageFile) TableName() string {
	return "messagefiles"
}
