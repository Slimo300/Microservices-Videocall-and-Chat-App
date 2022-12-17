package orm

import (
	"fmt"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DEFAULT_VERIFICATION_DURATION = 24 * time.Hour
const DEFAULT_RESET_DURATION = 15 * time.Minute

type Database struct {
	*gorm.DB
	Config DBConfig
}

type DBConfig struct {
	VerificationCodeDuration time.Duration
	ResetCodeDuration        time.Duration
}

type Option func(*Database)

func WithConfig(conf DBConfig) Option {
	return func(d *Database) {
		if conf.ResetCodeDuration != 0 {
			d.Config.ResetCodeDuration = conf.ResetCodeDuration
		}
		if conf.VerificationCodeDuration != 0 {
			d.Config.VerificationCodeDuration = conf.VerificationCodeDuration
		}
	}
}

// Setup creates Database object and initializes connection between MySQL database
func Setup(address string, options ...Option) (*Database, error) {
	conn, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?parseTime=true", address)), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	if err := conn.AutoMigrate(&models.User{}, models.VerificationCode{}, models.ResetCode{}); err != nil {
		return nil, err
	}

	db := &Database{DB: conn, Config: DBConfig{DEFAULT_VERIFICATION_DURATION, DEFAULT_RESET_DURATION}}

	for _, option := range options {
		option(db)
	}

	go db.CleanCodes(1 * time.Minute)

	return db, nil

}
