package redis

import (
	"fmt"
	"strings"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/google/uuid"
)

// NewMember sets new member in format <USER_ID>:<GROUP_ID>
func (db *DB) NewMember(event events.MemberCreatedEvent) error {
	return db.Set(fmt.Sprintf("%s:%s", event.UserID.String(), event.GroupID.String()), true, db.AccessCodeTTL).Err()
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

// GetUserGroups finds and returns all group ids that user is a member of
func (db *DB) GetUserGroups(userID uuid.UUID) ([]uuid.UUID, error) {
	keys, err := db.Keys(fmt.Sprintf("%s:*", userID.String())).Result()
	if err != nil {
		return nil, err
	}

	var groups []uuid.UUID
	for _, key := range keys {
		groupID := strings.Split(key, ":")[1]

		groupUID, err := uuid.Parse(groupID)
		if err != nil {
			return nil, err
		}

		groups = append(groups, groupUID)
	}

	return groups, nil
}
