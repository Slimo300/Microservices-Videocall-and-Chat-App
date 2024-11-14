package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	id       uuid.UUID
	groupID  uuid.UUID
	posted   time.Time
	text     string
	memberID uuid.UUID
	member   Member
	deleters []Member
	files    []MessageFile
}

func (m Message) ID() uuid.UUID        { return m.id }
func (m Message) GroupID() uuid.UUID   { return m.groupID }
func (m Message) MemberID() uuid.UUID  { return m.memberID }
func (m Message) Member() Member       { return m.member }
func (m Message) Text() string         { return m.text }
func (m Message) Posted() time.Time    { return m.posted }
func (m Message) Files() []MessageFile { return m.files }
func (m Message) Deleters() []Member   { return m.deleters }

func (m Message) IsUserDeleter(userID uuid.UUID) bool {
	for _, d := range m.deleters {
		if d.userID == userID {
			return true
		}
	}
	return false
}

func (m *Message) AddFiles(files ...MessageFile) {
	m.files = append(m.files, files...)
}

func (m *Message) AddDeleters(deleters ...Member) {
	m.deleters = append(m.deleters, deleters...)
}

func NewMessage(msgID, groupID, memberID uuid.UUID, text string, posted time.Time, files []MessageFile) Message {
	return Message{
		id:       msgID,
		groupID:  groupID,
		memberID: memberID,
		text:     text,
		posted:   posted,
		files:    files,
	}
}

func UnmarshalMessageFromDatabase(msgID, groupID, memberID uuid.UUID, text string, posted time.Time, member Member, deleters []Member, files []MessageFile) Message {
	return Message{
		id:       msgID,
		memberID: memberID,
		groupID:  groupID,
		text:     text,
		posted:   posted,
		member:   member,
		deleters: deleters,
		files:    files,
	}
}

type MessageFile struct {
	messageID uuid.UUID
	key       string
	extension string
}

func (m MessageFile) Key() string          { return m.key }
func (m MessageFile) MessageID() uuid.UUID { return m.messageID }
func (m MessageFile) Extension() string    { return m.extension }

func NewMessageFile(messageID uuid.UUID, key, ext string) MessageFile {
	return MessageFile{
		messageID: messageID,
		key:       key,
		extension: ext,
	}
}
