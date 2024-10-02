package orm

import (
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	UserName   string    `gorm:"column:username;unique"`
	Email      string    `gorm:"column:email;unique"`
	Password   string    `gorm:"column:password"`
	HasPicture *bool     `gorm:"column:has_picture"`
	Verified   *bool     `gorm:"column:verified"`
}

func (User) TableName() string {
	return "users"
}

func marshalUser(user models.User) User {
	hasPicture := user.HasPicture()
	verified := user.Verified()

	return User{
		ID:         user.ID(),
		UserName:   user.Username(),
		Email:      user.Email(),
		Password:   user.PasswordHash(),
		HasPicture: &hasPicture,
		Verified:   &verified,
	}
}

type AuthorizationCode struct {
	Code     uuid.UUID `gorm:"primaryKey"`
	UserID   uuid.UUID `gorm:"column:user_id;size:191"`
	Created  time.Time `gorm:"column:created"`
	CodeType string    `gorm:"column:code_type"`
}

func (AuthorizationCode) TableName() string {
	return "codes"
}

func marshalCode(code models.AuthorizationCode) AuthorizationCode {
	return AuthorizationCode{
		Code:     code.Code(),
		UserID:   code.UserID(),
		Created:  code.CreatedAt(),
		CodeType: code.Type().String(),
	}
}
