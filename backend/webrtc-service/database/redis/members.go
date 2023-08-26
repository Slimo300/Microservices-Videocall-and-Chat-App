package redis

import (
	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-service/models"
)

func (db *DB) GetMember(memberID string) (*models.Member, error) {
	member, err := db.HGetAll(memberID).Result()
	if err != nil {
		return nil, err
	}

	var muting bool
	if member["muting"] == "true" {
		muting = true
	}

	return &models.Member{
		ID:         memberID,
		GroupID:    member["groupID"],
		UserID:     member["userID"],
		Username:   member["username"],
		PictureURL: member["pictureURL"],
		Muting:     muting,
	}, nil
}

// NewMember sets new member in format <USER_ID>:<GROUP_ID>
func (db *DB) NewMember(member models.Member) error {

	pipe := db.TxPipeline()

	if err := pipe.SAdd(member.GroupID, member.ID).Err(); err != nil {
		return err
	}

	if err := pipe.HMSet(member.ID, map[string]interface{}{
		"groupID":    member.GroupID,
		"userID":     member.UserID,
		"username":   member.Username,
		"pictureURL": member.PictureURL,
		"muting":     member.Muting,
	}).Err(); err != nil {
		return err
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}

	return nil
}

// DeleteMember deletes given member if he exists
func (db *DB) DeleteMember(memberID string) error {
	return db.Del(memberID).Err()
}

// DeleteGroup deletes all members of a group
func (db *DB) DeleteGroup(groupID string) error {

	// Get all members of group
	members, err := db.SMembers(groupID).Result()
	if err != nil {
		return err
	}

	pipe := db.TxPipeline()

	if err := pipe.Del(members...).Err(); err != nil {
		return err
	}
	if err := pipe.Del(groupID).Err(); err != nil {
		return err
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}

	return nil
}
