package redis

import (
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
)

func (db *DB) GetMember(userID, groupID string) (string, error) {
	username, err := db.Get(fmt.Sprintf("%s:%s", userID, groupID)).Result()
	if err != nil {
		return "", err
	}

	return username, nil
}

// NewMember sets new member in format <USER_ID>:<GROUP_ID>
func (db *DB) NewMember(event events.MemberCreatedEvent) error {
	return db.Set(fmt.Sprintf("%s:%s", event.UserID.String(), event.GroupID.String()), event.User.UserName, 0).Err()
}

// DeleteMember deletes given member if he exists
func (db *DB) DeleteMember(event events.MemberDeletedEvent) error {
	return db.Del(fmt.Sprintf("%s:%s", event.UserID.String(), event.GroupID.String())).Err()
}

// DeleteGroup deletes all members of a group
func (db *DB) DeleteGroup(event events.GroupDeletedEvent) error {
	keys, err := db.Keys(fmt.Sprintf("*:%s", event.ID.String())).Result()
	if err != nil {
		return err
	}

	return db.Del(keys...).Err()
}
