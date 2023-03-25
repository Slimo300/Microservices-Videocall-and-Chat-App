package orm

import (
	"github.com/Slimo300/chat-messageservice/internal/models"
	"github.com/google/uuid"
)

func (db *Database) GetGroupMembership(userID, groupID uuid.UUID) (member models.Membership, err error) {
	return member, db.Where(models.Membership{UserID: userID, GroupID: groupID}).First(&member).Error
}
