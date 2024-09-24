package orm

import (
	"fmt"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/database"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GroupsGormRepository struct {
	db *gorm.DB
}

// Setup creates Database object and initializes connection between MySQL database
func Setup(dbaddress string) (*GroupsGormRepository, error) {

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?parseTime=true", dbaddress)), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	_ = db.AutoMigrate(&models.User{}, &models.Group{}, &models.Member{}, &models.Invite{})

	return &GroupsGormRepository{db: db}, nil
}

func (r *GroupsGormRepository) BeginTransaction() (database.GroupsRepository, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &GroupsGormRepository{db: tx}, nil
}

func (r *GroupsGormRepository) CommitTransaction() error {
	if err := r.db.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (r *GroupsGormRepository) RollbackTransaction() error {
	if err := r.db.Rollback().Error; err != nil {
		return err
	}
	return nil
}
