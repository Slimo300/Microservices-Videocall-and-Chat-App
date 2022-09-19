package models

import "github.com/google/uuid"

type VerificationCode struct {
	UserID         uuid.UUID `gorm:"column:user_id;size:191;primaryKey" json:"userID"`
	ActivationCode string    `gorm:"column:activation_code" json:"activation"`
}

func (VerificationCode) TableName() string {
	return "activation_codes"
}
