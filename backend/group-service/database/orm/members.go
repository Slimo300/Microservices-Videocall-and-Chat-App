package orm

import (
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"github.com/google/uuid"
)

func (db *Database) GetMemberByID(memberID uuid.UUID) (member *models.Member, err error) {
	return member, db.Preload("User").First(&member, memberID).Error
}

func (db *Database) GetMemberByUserGroupID(userID, groupID uuid.UUID) (member *models.Member, err error) {
	return member, db.Preload("User").Where(&models.Member{UserID: userID, GroupID: groupID}).First(&member).Error
}

func (db *Database) CreateMember(member *models.Member) (*models.Member, error) {
	if err := db.Create(&member).Error; err != nil {
		return nil, err
	}
	return member, db.Preload("User").First(&member, member.ID).Error
}

func (db *Database) UpdateMember(member *models.Member) (*models.Member, error) {
	return member, db.Preload("User").Save(&member).Error
}

func (db *Database) DeleteMember(memberID uuid.UUID) (member *models.Member, err error) {
	if err := db.Preload("User").First(&member, memberID).Error; err != nil {
		return nil, err
	}
	return member, db.Delete(&models.Member{}, memberID).Error
}
