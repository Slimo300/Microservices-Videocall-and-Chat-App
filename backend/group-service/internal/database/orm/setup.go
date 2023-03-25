package orm

import (
	"fmt"

	"github.com/Slimo300/chat-groupservice/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

// Setup creates Database object and initializes connection between MySQL database
func Setup(dbaddress string) (*Database, error) {

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?parseTime=true", dbaddress)), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Group{}, &models.Member{}, &models.Invite{}); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}
