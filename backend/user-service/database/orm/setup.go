package orm

import (
	"fmt"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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
func Setup(dbtype, address string, options ...Option) (*Database, error) {

	var dialector gorm.Dialector
	switch dbtype {
	case "MYSQL":
		dialector = mysql.Open(fmt.Sprintf("%s?parseTime=true", address))
	case "PostgreSQL":
		dialector = postgres.Open(address)
	default:
		return nil, fmt.Errorf("Unsupported database type: %s", dbtype)
	}

	conn, err := gorm.Open(dialector, &gorm.Config{
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

	go db.CleanCodes(time.Hour)

	return db, nil

}
