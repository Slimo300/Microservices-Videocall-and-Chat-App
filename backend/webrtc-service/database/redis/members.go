package redis

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/events"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/webrtc-service/models"
)

func StrToBool(str string) bool {
	if str == "0" || str == "" {
		return false
	} else {
		return true
	}
}

func (db *DB) GetMemberByID(memberID string) (*models.Member, error) {
	member, err := db.HGetAll(memberID).Result()
	if err != nil {
		return nil, err
	}

	return &models.Member{
		ID:       memberID,
		GroupID:  member["groupID"],
		UserID:   member["userID"],
		Username: member["username"],
		Muting:   StrToBool(member["muting"]),
		Admin:    StrToBool(member["admin"]),
		Creator:  StrToBool(member["creator"]),
	}, nil
}

func (db *DB) GetMemberByGroupAndUserID(groupID, userID string) (*models.Member, error) {
	memberID, err := db.HGet(groupID, userID).Result()
	if err != nil {
		return nil, err
	}

	member, err := db.HGetAll(memberID).Result()
	if err != nil {
		return nil, err
	}

	return &models.Member{
		ID:       memberID,
		GroupID:  member["groupID"],
		UserID:   member["userID"],
		Username: member["username"],
		Muting:   StrToBool(member["muting"]),
		Admin:    StrToBool(member["admin"]),
		Creator:  StrToBool(member["creator"]),
	}, nil
}

// NewMember sets new member in format <USER_ID>:<GROUP_ID>
func (db *DB) NewMember(member models.Member) error {

	pipe := db.TxPipeline()

	// In Redis we will store group object for searching for members by group and user id's
	// E.G.: GROUP_ID -> {
	//		USER_ID -> MEMBER_ID,
	// 		...
	// }
	if err := pipe.HSet(member.GroupID, member.UserID, member.ID).Err(); err != nil {
		return err
	}

	// In Redis we will store member object containing all member data
	if err := pipe.HMSet(member.ID, map[string]interface{}{
		"groupID":  member.GroupID,
		"userID":   member.UserID,
		"username": member.Username,
		"muting":   member.Muting,
		"admin":    member.Admin,
		"creator":  member.Creator,
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

	member, err := db.HGetAll(memberID).Result()
	if err != nil {
		return err
	}

	if err := db.HDel(member["groupID"], member["userID"]).Err(); err != nil {
		return err
	}

	return db.Del(memberID).Err()
}

// DeleteGroup deletes all members of a group
func (db *DB) DeleteGroup(groupID string) error {

	// Get all members of group
	members, err := db.HGetAll(groupID).Result()
	if err != nil {
		return err
	}

	pipe := db.TxPipeline()

	for _, memberID := range members {
		if err := pipe.Del(memberID).Err(); err != nil {
			return err
		}
	}
	if err := pipe.Del(groupID).Err(); err != nil {
		return err
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}

	return nil
}

func (db *DB) ModifyMember(event events.MemberUpdatedEvent) error {
	return db.HSet(event.ID.String(), "muting", event.Muting).Err()
}
