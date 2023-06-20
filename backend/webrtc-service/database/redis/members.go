package redis

import (
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

func (db *DB) IsUserMember(userID, groupID string) bool {
	_, err := db.Keys(fmt.Sprintf("%s:%s", userID, groupID)).Result()
	if err != nil {
		return false
	}

	return true
}

// NewMember sets new member in format <USER_ID>:<GROUP_ID>
func (db *DB) NewMember(event events.MemberCreatedEvent) error {
	return db.Set(fmt.Sprintf("%s:%s", event.UserID.String(), event.GroupID.String()), true, 0).Err()
}

// DeleteMember deletes given member if he exists
func (db *DB) DeleteMember(event events.MemberDeletedEvent) error {
	return db.Del(fmt.Sprintf("%s:%s", event.UserID.String(), event.GroupID.String())).Err()
}

// DeleteGroup deletes all members of a group
func (db *DB) DeleteGroup(event events.GroupDeletedEvent) error {
	res := db.Keys(fmt.Sprintf("*:%s", event.ID.String()))
	if err := res.Err(); err != nil {
		return err
	}
	keys := res.Val()

	return db.Del(keys...).Err()
}
