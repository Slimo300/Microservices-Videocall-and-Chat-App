package orm

import (
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	UserName   string    `gorm:"unique"`
	Email      string    `gorm:"uniqueIndex,type:VARCHAR(255)"`
	Password   string
	HasPicture *bool
	Verified   *bool
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
	Code     uuid.UUID `gorm:"primaryKey,priority=2"`
	CodeType string    `gorm:"primaryKey,priority=1"`
	UserID   uuid.UUID `gorm:"size:191"`
	User     User      `gorm:"constraint:OnDelete:CASCADE"`
	Created  time.Time
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
