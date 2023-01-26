package orm

import (
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

// Setup creates Database object and initializes connection between MySQL database
func Setup(dbtype, address string) (*Database, error) {
	var dialector gorm.Dialector
	switch dbtype {
	case "MYSQL":
		dialector = mysql.Open(fmt.Sprintf("%s?parseTime=true", address))
	case "PostgreSQL":
		dialector = postgres.Open(address)
	default:
		return nil, fmt.Errorf("Unsupported database type: %s", dbtype)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Message{}, &models.Membership{}, &models.MessageFile{}); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}
