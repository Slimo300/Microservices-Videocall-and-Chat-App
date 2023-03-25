package models

import (
	"time"

	"github.com/google/uuid"
)

type ResetCode struct {
	UserID    uuid.UUID `gorm:"column:user_id;size:191;primaryKey" json:"userID"`
	ResetCode string    `gorm:"column:reset_code" json:"activation"`
	Created   time.Time `gorm:"column:created" json:"created"`
}

func (ResetCode) TableName() string {
	return "reset_codes"
}
