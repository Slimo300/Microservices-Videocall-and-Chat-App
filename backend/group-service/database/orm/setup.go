package orm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GroupsGormRepository struct {
	db *gorm.DB
}

func NewGroupsGormRepository(dbaddress string) (*GroupsGormRepository, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?parseTime=true", dbaddress)), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(User{}, Group{}, Member{}, Invite{}); err != nil {
		return nil, err
	}
	return &GroupsGormRepository{db: db}, nil
}
