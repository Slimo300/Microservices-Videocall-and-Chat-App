package database

import (
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/ws-service/models"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func Setup(address string) (*Database, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?parseTime=true", address)), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Membership{})

	return &Database{DB: db}, nil
}

func (db *Database) GetUserGroups(userID uuid.UUID) (groups []uuid.UUID, err error) {
	return groups, db.Where(models.Membership{UserID: userID}).Select("memberships.group_id").Scan(&groups).Error
}
