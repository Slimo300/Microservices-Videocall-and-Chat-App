package orm

import (
	"fmt"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

// Setup creates Database object and initializes connection between MySQL database
func Setup(dbaddress string) (*Database, error) {

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?parseTime=true", dbaddress)), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	_ = db.AutoMigrate(&models.User{}, &models.Group{}, &models.Member{}, &models.Invite{})

	return &Database{DB: db}, nil
}
