package orm

import (
	"fmt"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MessagesGormRepository struct {
	*gorm.DB
}

func NewMessagesGormRepository(address string) (database.MessagesRepository, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?parseTime=true", address)), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	_ = db.AutoMigrate(&Message{}, &Member{}, &MessageFile{})

	return &MessagesGormRepository{DB: db}, nil
}
