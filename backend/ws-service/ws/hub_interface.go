package ws

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

type Hub interface {
	Run()
	Join(*client)
	Leave(*client)
	Forward(*Message)

	GroupDeleted(event events.GroupDeletedEvent)
	MemberAdded(event events.MemberCreatedEvent)
	MemberDeleted(event events.MemberDeletedEvent)
}
